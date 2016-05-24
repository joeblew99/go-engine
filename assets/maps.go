package assets

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/walesey/go-engine/editor/models"
	"github.com/walesey/go-engine/renderer"
)

func LoadMap(path string) *editorModels.MapModel {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("Error Reading map file: %v\n", err)
		return nil
	}

	var mapModel editorModels.MapModel
	err = json.Unmarshal(data, &mapModel)
	if err != nil {
		log.Printf("Error unmarshaling map model: %v\n", err)
		return nil
	}

	return &mapModel
}

func LoadMapToNode(srcModel *editorModels.NodeModel, destNode *renderer.Node) *editorModels.NodeModel {
	copy := srcModel.Copy(func(name string) string { return name })
	loadMapRecursive(copy, srcModel, destNode)
	return copy
}

func loadMapRecursive(model, srcModel *editorModels.NodeModel, destNode *renderer.Node) {
	model.SetNode(destNode)
	if model.Geometry != nil {
		geometry, err := ImportObjCached(*model.Geometry)
		if err == nil {
			destNode.Add(geometry)
		}
	}
	destNode.SetScale(model.Scale)
	destNode.SetTranslation(model.Translation)
	destNode.SetOrientation(model.Orientation)
	if model.Reference != nil {
		if refModel := FindNodeById(*model.Reference, srcModel); refModel != nil {
			for _, childModel := range refModel.Children {
				model.Children = append(model.Children, childModel.Copy(func(name string) string { //TODO: fix this for refs within refs
					return fmt.Sprintf("%v::%v", *model.Reference, name)
				}))
			}
		}
	}
	for _, childModel := range model.Children {
		newNode := renderer.CreateNode()
		destNode.Add(newNode)
		loadMapRecursive(childModel, srcModel, newNode)
	}
}

func FindNodeById(nodeId string, model *editorModels.NodeModel) *editorModels.NodeModel {
	if model.Id == nodeId {
		return model
	}
	for _, childModel := range model.Children {
		node := FindNodeById(nodeId, childModel)
		if node != nil {
			return node
		}
	}
	return nil
}

func FindNodeByClass(class string, model *editorModels.NodeModel) []*editorModels.NodeModel {
	results := []*editorModels.NodeModel{}
	if hasClass(class, model) {
		results = append(results, model)
	}
	for _, childModel := range model.Children {
		results = append(results, FindNodeByClass(class, childModel)...)
	}
	return results
}

func hasClass(class string, model *editorModels.NodeModel) bool {
	for _, modelClass := range model.Classes {
		if class == modelClass {
			return true
		}
	}
	return false
}
