package file

import (
	"bytes"

	"github.com/shrutikamendhe/dockm/api"

	"io"
	"os"
	"path"
	"strconv"
)

const (
	// TLSStorePath represents the subfolder where TLS files are stored in the file store folder.
	TLSStorePath = "tls"
	// TLSCACertFile represents the name on disk for a TLS CA file.
	TLSCACertFile = "ca.pem"
	// TLSCertFile represents the name on disk for a TLS certificate file.
	TLSCertFile = "cert.pem"
	// TLSKeyFile represents the name on disk for a TLS key file.
	TLSKeyFile = "key.pem"
	// ComposeStorePath represents the subfolder where compose files are stored in the file store folder.
	ComposeStorePath = "compose"
)

// Service represents a service for managing files and directories.
type Service struct {
	dataStorePath string
	fileStorePath string
}

// NewService initializes a new service. It creates a data directory and a directory to store files
// inside this directory if they don't exist.
func NewService(dataStorePath, fileStorePath string) (*Service, error) {
	service := &Service{
		dataStorePath: dataStorePath,
		fileStorePath: path.Join(dataStorePath, fileStorePath),
	}

	// Checking if a mount directory exists is broken with Go on Windows.
	// This will need to be reviewed after the issue has been fixed in Go.
	// See: https://github.com/click2cloud/dockm/issues/474
	// err := createDirectoryIfNotExist(dataStorePath, 0755)
	// if err != nil {
	// 	return nil, err
	// }

	err := service.createDirectoryInStoreIfNotExist(TLSStorePath)
	if err != nil {
		return nil, err
	}

	err = service.createDirectoryInStoreIfNotExist(ComposeStorePath)
	if err != nil {
		return nil, err
	}

	return service, nil
}

// StoreComposeEnvFile stores a new .env file in the stack store path using the content of envFileContent.
func (service *Service) StoreComposeEnvFile(name, envFileContent string) error {
	stackStorePath := path.Join(ComposeStorePath, name)
	err := service.createDirectoryInStoreIfNotExist(stackStorePath)
	if err != nil {
		return err
	}

	envFilePath := path.Join(stackStorePath, ".env")
	data := []byte(envFileContent)
	r := bytes.NewReader(data)

	err = service.createFileInStore(envFilePath, r)
	if err != nil {
		return err
	}

	return nil
}

// StoreComposeFile creates a subfolder in the ComposeStorePath and stores a new file using the content from composeFileContent.
// It returns the path to the folder where the file is stored.
func (service *Service) StoreComposeFile(name, composeFileContent string) (string, error) {
	stackStorePath := path.Join(ComposeStorePath, name)
	err := service.createDirectoryInStoreIfNotExist(stackStorePath)
	if err != nil {
		return "", err
	}

	composeFilePath := path.Join(stackStorePath, "docker-compose.yml")
	data := []byte(composeFileContent)
	r := bytes.NewReader(data)

	err = service.createFileInStore(composeFilePath, r)
	if err != nil {
		return "", err
	}

	return path.Join(service.fileStorePath, stackStorePath), nil
}

// StoreTLSFile creates a subfolder in the TLSStorePath and stores a new file with the content from r.
func (service *Service) StoreTLSFile(endpointID dockm.EndpointID, fileType dockm.TLSFileType, r io.Reader) error {
	ID := strconv.Itoa(int(endpointID))
	endpointStorePath := path.Join(TLSStorePath, ID)
	err := service.createDirectoryInStoreIfNotExist(endpointStorePath)
	if err != nil {
		return err
	}

	var fileName string
	switch fileType {
	case dockm.TLSFileCA:
		fileName = TLSCACertFile
	case dockm.TLSFileCert:
		fileName = TLSCertFile
	case dockm.TLSFileKey:
		fileName = TLSKeyFile
	default:
		return dockm.ErrUndefinedTLSFileType
	}

	tlsFilePath := path.Join(endpointStorePath, fileName)
	err = service.createFileInStore(tlsFilePath, r)
	if err != nil {
		return err
	}
	return nil
}

// GetPathForTLSFile returns the absolute path to a specific TLS file for an endpoint.
func (service *Service) GetPathForTLSFile(endpointID dockm.EndpointID, fileType dockm.TLSFileType) (string, error) {
	var fileName string
	switch fileType {
	case dockm.TLSFileCA:
		fileName = TLSCACertFile
	case dockm.TLSFileCert:
		fileName = TLSCertFile
	case dockm.TLSFileKey:
		fileName = TLSKeyFile
	default:
		return "", dockm.ErrUndefinedTLSFileType
	}
	ID := strconv.Itoa(int(endpointID))
	return path.Join(service.fileStorePath, TLSStorePath, ID, fileName), nil
}

// DeleteStackFiles deletes a folder containing all the files associated to a stack.
func (service *Service) DeleteStackFiles(projectPath string) error {
	err := os.RemoveAll(projectPath)
	if err != nil {
		return err
	}
	return nil
}

// DeleteTLSFiles deletes a folder containing the TLS files for an endpoint.
func (service *Service) DeleteTLSFiles(endpointID dockm.EndpointID) error {
	ID := strconv.Itoa(int(endpointID))
	endpointPath := path.Join(service.fileStorePath, TLSStorePath, ID)
	err := os.RemoveAll(endpointPath)
	if err != nil {
		return err
	}
	return nil
}

// createDirectoryInStoreIfNotExist creates a new directory in the file store if it doesn't exists on the file system.
func (service *Service) createDirectoryInStoreIfNotExist(name string) error {
	path := path.Join(service.fileStorePath, name)
	return createDirectoryIfNotExist(path, 0700)
}

// createDirectoryIfNotExist creates a directory if it doesn't exists on the file system.
func createDirectoryIfNotExist(path string, mode uint32) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, os.FileMode(mode))
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

// createFile creates a new file in the file store with the content from r.
func (service *Service) createFileInStore(filePath string, r io.Reader) error {
	path := path.Join(service.fileStorePath, filePath)
	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, r)
	if err != nil {
		return err
	}
	return nil
}
