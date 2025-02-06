package cmd

import (
	"maps"
	"net/url"
	"path/filepath"
	"regexp"

	"github.com/apple/pkl-go/pkl"
	"github.com/spf13/cobra"
	"hpkl.io/hpkl/pkg/app"
	"hpkl.io/hpkl/pkg/pklutils"
)

func NewResolveCmd(appConfig *app.AppConfig) *cobra.Command {

	logger := appConfig.Logger

	cmd := &cobra.Command{
		Use:   "resolve",
		Short: "Resolve all dependencies from pkl project",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				for _, v := range args {
					logger.Info("Resolving path: %s", v)
					appConfig.WorkingDir = v
					appConfig.Reset()
					err := Resolve(appConfig)
					if err != nil {
						panic(err)
					}
				}
			} else {
				err := Resolve(appConfig)
				if err != nil {
					panic(err)
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&appConfig.PlainHttp, "plain-http", "p", false, "Use plain http for registry")

	return cmd
}

func CollectLocalDependencies(dependecies *pkl.ProjectDependencies) map[string]*app.Dependency {
	result := make(map[string]*app.Dependency)

	for n, dep := range dependecies.LocalDependencies {
		result[dep.PackageUri] = &app.Dependency{Uri: dep.PackageUri, Name: n, ProjectFileUri: dep.ProjectFileUri}

		for _, localDep := range dep.Dependencies.LocalDependencies {
			inner := CollectLocalDependencies(localDep.Dependencies)
			maps.Copy(result, inner)
		}
	}

	return result
}

func CollectRemoteDependencies(dependecies *pkl.ProjectDependencies) map[string]app.Dependency {
	result := make(map[string]app.Dependency)

	for _, dep := range dependecies.LocalDependencies {
		remote := dep.Dependencies.RemoteDependencies

		for n, remoteDep := range remote {
			result[remoteDep.PackageUri] = app.Dependency{Uri: remoteDep.PackageUri, Name: n}
		}

		for _, localDep := range dep.Dependencies.LocalDependencies {
			inner := CollectRemoteDependencies(localDep.Dependencies)
			maps.Copy(result, inner)
		}
	}

	for n, dep := range dependecies.RemoteDependencies {
		result[dep.PackageUri] = app.Dependency{Uri: dep.PackageUri, Name: n}
	}

	return result
}

func Resolve(appConfig *app.AppConfig) error {
	resolver, err := app.NewResolver(appConfig)
	if err != nil {
		appConfig.Logger.Error("Error on creating resolver")
		return err
	}

	project := appConfig.Project()

	remoteDependencies := CollectRemoteDependencies(project.Dependencies())

	resolvedDependencies, err := resolver.Resolve(remoteDependencies)

	if err != nil {
		appConfig.Logger.Error("Error on resolving remote dependencies")
		return err
	}

	resolvedDependencies, err = resolver.Deduplicate(resolvedDependencies)

	if err != nil {
		appConfig.Logger.Error("Error on deduplication")
		return err
	}

	err = resolver.Download(resolvedDependencies)

	if err != nil {
		return err
	}

	dependencies := make(map[string]*pklutils.ResolvedDependency, len(resolvedDependencies)+len(project.Dependencies().LocalDependencies))

	projectDeps := pklutils.ProjectDeps{
		SchemaVersion:        1,
		ResolvedDependencies: dependencies,
	}

	for _, dep := range resolvedDependencies {

		mapUri, err := resolver.MajorVersionPackage(dep)

		if err != nil {
			appConfig.Logger.Error("Error on dependency resolving")
			return err
		}

		packageUri, err := url.Parse(dep.PackageUri)

		if err != nil {
			return err
		}

		packageUri.Scheme = "projectpackage"

		resolvedDependency := pklutils.ResolvedDependency{
			DependencyType: "remote",
			Uri:            packageUri.String(),
			Checksums:      map[string]string{"sha256": dep.Checksum},
		}

		projectDeps.ResolvedDependencies[mapUri] = &resolvedDependency
	}

	projectFileUri, err := url.Parse(project.ProjectFileUri)
	if err != nil {
		appConfig.Logger.Error("Error on Url Parsing")
		return err
	}

	projectFilePath := filepath.Dir(projectFileUri.Path)

	versionRegex := regexp.MustCompile(`^(.*)@(\d+)`)

	localDependencies := CollectLocalDependencies(project.Dependencies())

	for _, dep := range localDependencies {

		projectUri, err := url.Parse(dep.Uri)
		if err != nil {
			appConfig.Logger.Error("Error on Url Parsing in dependency")
			return err
		}
		projectUri.Scheme = "projectpackage"

		mapUri := versionRegex.FindStringSubmatch(dep.Uri)[0]

		depProjectFileUri, err := url.Parse(dep.ProjectFileUri)
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(projectFilePath, filepath.Dir(depProjectFileUri.Path))

		if err != nil {
			appConfig.Logger.Error("Error on Url Parsing in dependency")
			return err
		}

		resolvedDependency := pklutils.ResolvedDependency{
			DependencyType: "local",
			Path:           rel,
			Uri:            projectUri.String(),
		}
		projectDeps.ResolvedDependencies[mapUri] = &resolvedDependency
	}

	err = pklutils.PklWriteDeps(appConfig.WorkingDir, &projectDeps)
	if err != nil {
		appConfig.Logger.Error("Error on write deps")
		return err
	}

	return nil
}
