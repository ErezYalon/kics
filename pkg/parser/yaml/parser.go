package json

import (
	"bytes"
	"encoding/json"

	"github.com/Checkmarx/kics/pkg/model"
	"github.com/Checkmarx/kics/pkg/parser/utils"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Parser defines a parser type
type Parser struct {
}

// Playbooks represents a playbook object from parsed yaml files
type Playbooks struct {
	Tasks []map[string]interface{} `json:"playbooks"`
}

// Resolve - replace or modifies in-memory content before parsing
func (p *Parser) Resolve(fileContent []byte, filename string) (*[]byte, error) {
	return &fileContent, nil
}

// Parse parses yaml/yml file and returns it as a Document
func (p *Parser) Parse(filePath string, fileContent []byte) ([]model.Document, error) {
	var documents []model.Document
	dec := yaml.NewDecoder(bytes.NewReader(fileContent))

	doc := &model.Document{}
	for dec.Decode(doc) == nil {
		if doc != nil {
			documents = append(documents, *doc)
		}
		doc = &model.Document{}
	}

	if documents == nil {
		var err error
		documents, err = playbookParser(filePath, fileContent)
		if err != nil {
			return nil, errors.Wrap(err, "failed to Parse YAML")
		}
		if documents == nil {
			return nil, errors.Wrap(errors.New("invalid yaml"), "failed to parse yaml")
		}
	}

	return convertKeysToString(documents), nil
}

// convertKeysToString goes through every document to convert map[interface{}]interface{}
// to map[string]interface{}
func convertKeysToString(docs []model.Document) []model.Document {
	documents := make([]model.Document, 0, len(docs))
	for _, doc := range docs {
		for key, value := range doc {
			doc[key] = convert(value)
		}
		documents = append(documents, doc)
	}
	return documents
}

// convert goes recursively through the keys in the given value and converts nested maps type of map[interface{}]interface{}
// to map[string]interface{}
func convert(value interface{}) interface{} {
	switch t := value.(type) {
	case map[interface{}]interface{}:
		mapStr := map[string]interface{}{}
		for key, val := range t {
			if t, ok := key.(string); ok {
				mapStr[t] = convert(val)
			}
		}
		return mapStr
	case []interface{}:
		for key, val := range t {
			t[key] = convert(val)
		}
	case model.Document:
		for key, val := range t {
			t[key] = convert(val)
		}
	}
	return value
}

// SupportedExtensions returns extensions supported by this parser, which are yaml and yml extension
func (p *Parser) SupportedExtensions() []string {
	return []string{".yaml", ".yml"}
}

// SupportedTypes returns types supported by this parser, which are ansible, cloudFormation, k8s
func (p *Parser) SupportedTypes() []string {
	return []string{"Ansible", "CloudFormation", "Kubernetes", "OpenAPI"}
}

// GetKind returns YAML constant kind
func (p *Parser) GetKind() model.FileKind {
	return model.KindYAML
}

func processSwaggerContent(elements map[string]interface{}, filePath string) {
	swaggerInfo := utils.AddSwaggerInfo(filePath, elements["swagger_file"].(string))
	if swaggerInfo != nil {
		elements["swagger_file"] = swaggerInfo
	}
}

func processCertContent(elements map[string]interface{}, content, filePath string) {
	var certInfo map[string]interface{}
	if content != "" {
		certInfo = utils.AddCertificateInfo(filePath, content)
		if certInfo != nil {
			elements["certificate"] = certInfo
		}
	}
}

func processElements(elements map[string]interface{}, filePath string) {
	if elements["certificate"] != nil {
		processCertContent(elements, utils.CheckCertificate(elements["certificate"].(string)), filePath)
	}
	if elements["swagger_file"] != nil {
		processSwaggerContent(elements, filePath)
	}
}

func addExtraInfo(documents []model.Document, filePath string) []model.Document {
	for _, documentPlaybooks := range documents { // iterate over documents
		for _, resources := range documentPlaybooks["playbooks"].([]interface{}) { // iterate over playbooks
			for _, v := range resources.(map[string]interface{}) {
				_, ok := v.(map[string]interface{})
				if ok {
					processElements(v.(map[string]interface{}), filePath)
				}
			}
		}
	}

	return documents
}

func playbookParser(filePath string, fileContent []byte) ([]model.Document, error) {
	doc := &model.Document{}
	dec := yaml.NewDecoder(bytes.NewReader(fileContent))
	arr := make([]map[string]interface{}, 0)
	var playBooks Playbooks
	var documents []model.Document
	for dec.Decode(&arr) == nil {
		if doc != nil {
			playBooks.Tasks = append(playBooks.Tasks, arr...)
			j, err := json.Marshal(playBooks)
			if err != nil {
				return nil, errors.Wrap(err, "failed to Marshal YAML")
			}

			if err := json.Unmarshal(j, &doc); err != nil {
				return nil, errors.Wrap(err, "failed to Unmarshal YAML")
			}
			documents = append(documents, *doc)
		}
		doc = &model.Document{}
		arr = make([]map[string]interface{}, 0)
	}

	return addExtraInfo(documents, filePath), nil
}

// GetCommentToken return the comment token of YAML - #
func (p *Parser) GetCommentToken() string {
	return "#"
}
