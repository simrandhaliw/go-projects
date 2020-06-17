package index

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/util/templates"

	"github.com/operator-framework/operator-registry/pkg/containertools"
	"github.com/operator-framework/operator-registry/pkg/lib/indexer"
)

var exportLong = templates.LongDesc(`
	Export an operator from an index image into the appregistry format. 

	This command will take an index image (specified by the --index option), parse it for the given operator (set by 
	the --operator option) and export the operator metadata into an appregistry compliant format (a package.yaml file). 
	This command requires access to docker or podman to complete successfully.

	Note: the appregistry format is being deprecated in favor of the new index image and image bundle format. 
	`)

func newIndexExportCmd() *cobra.Command {
	indexCmd := &cobra.Command{
		Use:   "export",
		Short: "Export an operator from an index into the appregistry format",
		Long:  exportLong,

		PreRunE: func(cmd *cobra.Command, args []string) error {
			if debug, _ := cmd.Flags().GetBool("debug"); debug {
				logrus.SetLevel(logrus.DebugLevel)
			}
			return nil
		},

		RunE: runIndexExportCmdFunc,
	}

	indexCmd.Flags().Bool("debug", false, "enable debug logging")
	indexCmd.Flags().StringP("index", "i", "", "index to get package from")
	if err := indexCmd.MarkFlagRequired("index"); err != nil {
		logrus.Panic("Failed to set required `index` flag for `index export`")
	}
	indexCmd.Flags().StringP("package", "o", "", "the package to export")
	if err := indexCmd.MarkFlagRequired("package"); err != nil {
		logrus.Panic("Failed to set required `package` flag for `index export`")
	}
	indexCmd.Flags().StringP("download-folder", "f", "downloaded", "directory where downloaded operator bundle(s) will be stored")
	indexCmd.Flags().StringP("container-tool", "c", "podman", "tool to interact with container images (save, build, etc.). One of: [none, docker, podman]")
	if err := indexCmd.Flags().MarkHidden("debug"); err != nil {
		logrus.Panic(err.Error())
	}

	return indexCmd

}

func runIndexExportCmdFunc(cmd *cobra.Command, args []string) error {
	index, err := cmd.Flags().GetString("index")
	if err != nil {
		return err
	}

	packageName, err := cmd.Flags().GetString("package")
	if err != nil {
		return err
	}

	downloadPath, err := cmd.Flags().GetString("download-folder")
	if err != nil {
		return err
	}

	containerTool, err := cmd.Flags().GetString("container-tool")
	if err != nil {
		return err
	}

	if containerTool == "none" {
		return fmt.Errorf("none is not a valid container-tool for index add")
	}

	logger := logrus.WithFields(logrus.Fields{"index": index, "package": packageName})

	logger.Info("export from the index")

	indexExporter := indexer.NewIndexExporter(containertools.NewContainerTool(containerTool, containertools.PodmanTool), logger)

	request := indexer.ExportFromIndexRequest{
		Index:         index,
		Package:       packageName,
		DownloadPath:  downloadPath,
		ContainerTool: containerTool,
	}

	err = indexExporter.ExportFromIndex(request)
	if err != nil {
		return err
	}

	return nil
}
