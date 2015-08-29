package assets

import (
	"strings"
	"strconv"
	"bufio"
	"os"
	"log"
	"image"

	"github.com/go-gl/mathgl/mgl32"
)

//vericies format : x,y,z,   nx,ny,nz,tx,ty,tz,btx,bty,btz,   u,v
//Indicies format : f1,f2,f3 (triangles)
type ObjData struct {
	Name string
	Indicies []uint32
	Vertices []float32
	Mtl *mtlData
}

type mtlData struct {
	Name string
	Ns,Ka,Kd,Ks,Ni,D float32
	Illum int
	Map_Kd, Map_Disp, Map_Spec, Map_Roughness image.Image
}

//returns corresponding index (0,1,2...)
func (obj *ObjData) pushVert( x,y,z,nx,ny,nz,tx,ty,tz,btx,bty,btz,u,v float32 ) uint32 {
	obj.Vertices = append(obj.Vertices, x,y,z,nx,ny,nz,tx,ty,tz,btx,bty,btz,u,v )
	return (uint32)((len(obj.Vertices) / 14) - 1)
}

func (obj *ObjData) pushIndex( indicies ...uint32 ) {
	obj.Indicies = append(obj.Indicies, indicies... )
}

//parses a single triangle vertex, returning the newly generated index
func (obj *ObjData) processFaceVertex( token string, vertexList, uvList, normalList []float32 ) uint32 {
	face := strings.Split(token, "/")
	var index int32

	//vertex
	vx := (float32)(0.0)
	vy := (float32)(0.0)
	vz := (float32)(0.0)
	if len(face) > 0 && face[0] != "" {
		index = (sti(face[0])-1) * 3
		vx = vertexList[index]
		vy = vertexList[index+1]
		vz = vertexList[index+2]
	}

	//texture
	vtx := (float32)(0.0)
	vty := (float32)(0.0)
	if len(face) > 1 && face[1] != "" {
		index = (sti(face[1])-1) * 2
		vtx = uvList[index]
		vty = uvList[index+1]
	}

	//normal / tangents
	nx := (float32)(0.0)
	ny := (float32)(0.0)
	nz := (float32)(0.0)
	tanx := (float32)(0.0)
	tany := (float32)(0.0)
	tanz := (float32)(0.0)
	bitanx := (float32)(0.0)
	bitany := (float32)(0.0)
	bitanz := (float32)(0.0)
	if len(face) > 2 && face[2] != "" {
		index = (sti(face[2])-1) * 3
		nx = normalList[index]
		ny = normalList[index+1]
		nz = normalList[index+2]
		norm := mgl32.Vec3{nx,ny,nz}
		tangent := norm.Cross(norm.Add(mgl32.Vec3{1,1,1}))
		bitangent := norm.Cross(tangent)
		tanx = tangent.X()
		tany = tangent.Y()
		tanz = tangent.Z()
		bitanx = bitangent.X()
		bitany = bitangent.Y()
		bitanz = bitangent.Z()
	}

	return obj.pushVert( vx,vy,vz, nx,ny,nz, tanx,tany,tanz, bitanx,bitany,bitanz, vtx,vty )
}

//Processes a polygonal face by splitting it into triangles
func (obj *ObjData) processFace( line string, vertexList, uvList, normalList []float32 ){
	tokens := strings.Fields(line)
	if tokens[0] == "f" {
		tokens = append(tokens[:0], tokens[1:]...)
		for len(tokens) > 0 {
			tempTokens := make([]string, 0, 0)
			for i:=0 ; i < (len(tokens)-1); i+=2 {
				obj.pushIndex( obj.processFaceVertex(tokens[i], vertexList, uvList, normalList) )
				obj.pushIndex( obj.processFaceVertex(tokens[i+1], vertexList, uvList, normalList) )
				if len(tokens) > (i+2){
					obj.pushIndex( obj.processFaceVertex(tokens[i+2], vertexList, uvList, normalList) )
				} else {
					obj.pushIndex( obj.processFaceVertex(tokens[0], vertexList, uvList, normalList) )
				}
				if len(tokens) > 4 {
					tempTokens = append(tempTokens, tokens[i])
				}
			}
			if len(tokens) > 4 && len(tokens)%2 == 1 {
				tempTokens = append(tempTokens, tokens[len(tokens)-1])
			}
			tokens = tempTokens
		}
	}
}

//imports an obj filePath into an ObjData reference containing index and vertex buffers
func ImportObj(filePath string) (*ObjData, error) {

	obj := &ObjData{ Indicies: make([]uint32, 0, 0), Vertices: make([]float32, 0, 0) }
	vertexList := make([]float32, 0, 0)
	uvList := make([]float32, 0, 0)
	normalList := make([]float32, 0, 0)

	//split the file name from the file path
	filePathTokens := strings.Split(filePath, "/")
	fileName := filePathTokens[ len(filePathTokens)-1 ]
	path := strings.TrimRight(filePath, fileName)

	//open the file and read all lines
	file, err := os.Open(filePath)
	if err != nil {
		return obj, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Fields(line)
		if len(tokens) > 0 {
			dataType := tokens[0]
			if dataType == "o" { //sub mesh
				obj.Name = tokens[1]
			} else if dataType == "v" { //xyz vertex
				vertexList = append(vertexList, stf(tokens[1]), stf(tokens[2]), stf(tokens[3]) )
			} else if dataType == "vt" { //uv coord
				uvList = append(uvList, stf(tokens[1]), stf(tokens[2]) )
			} else if dataType == "vn" { //xyz vertex normal
				normalList = append(normalList, stf(tokens[1]), stf(tokens[2]), stf(tokens[3]) )
			} else if dataType == "f" { // v/t/n face
				obj.processFace( line, vertexList, uvList, normalList );
			} else if dataType == "mtllib" {
				mtl,err := importMTL( path, tokens[1] )
				if err != nil {
					return obj, err
				}
				obj.Mtl = mtl
			} else if dataType == "usemtl" { //mtl material
				//TODO: multiple mtls
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return obj, err
	}

	return obj, nil
}

//Returns mtl object data type
func importMTL( filePath, fileName string ) (*mtlData, error){
 	mtl := &mtlData{}

	file, err := os.Open(filePath + fileName)
	if err != nil {
		return mtl, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Fields(line)
		if len(tokens) > 0 {
			dataType := tokens[0]
			if dataType == "newmtl" {
				mtl.Name = tokens[1]
			} else if dataType == "Ns" {
				mtl.Ns = stf(tokens[1])
				//TODO: Other mtl variables
			} else if dataType == "map_Kd" {
				mtl.Map_Kd = ImportImage( filePath + tokens[1] )
			} else if dataType == "map_Disp" {
				mtl.Map_Disp = ImportImage( filePath + tokens[1] )
			} else if dataType == "map_Spec" {
				mtl.Map_Spec = ImportImage( filePath + tokens[1] )
			} else if dataType == "map_Roughness" {
				mtl.Map_Roughness = ImportImage( filePath + tokens[1] )
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return mtl, err
	}

	return mtl, nil
}

//string to float32
func stf(s string) float32 {
	f, err := strconv.ParseFloat(s, 32)
	if err != nil {
	    log.Fatal(err)
	}
	return (float32)(f)
}

//string to int32
func sti(s string) int32 {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
	    log.Fatal(err)
	}
	return (int32)(i)
}