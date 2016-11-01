package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Parser struct {
	in   *bufio.Reader
	path string

	fragOut, vertOut, geoOut io.Writer

	token   Token  // The token
	literal string // The literal of the token, if any
	chr     rune

	frag, vert, geo bool
}

func New(src io.Reader, path string, vert, frag, geo io.Writer) *Parser {
	return &Parser{
		in:      bufio.NewReader(src),
		path:    path,
		fragOut: frag,
		vertOut: vert,
		geoOut:  geo,
	}
}

func ParseFile(path string, vert, frag, geo io.Writer) error {
	src, err := os.Open(path)
	if err != nil {
		return err
	}
	defer src.Close()

	New(src, path, vert, frag, geo).Parse()
	return nil
}

func (p *Parser) Parse() {
	for {
		p.next()
		if p.token == EOF {
			break
		}
		p.parseToken()
	}
}

func (p *Parser) parseToken() {
	switch p.token {
	case HASH:
		p.parseHash()
	case STRING:
		p.parseString()
	default:
		p.write(p.literal)
	}
}

func (p *Parser) parseHash() {
	literal := p.literal
	p.next()
	switch p.token {
	case INCLUDE:
		p.parseHashInclude()
	case VERT:
		p.parseHashVert()
	case FRAG:
		p.parseHashFrag()
	case GEO:
		p.parseHashGeo()
	case ENDVERT:
		p.parseHashEndVert()
	case ENDFRAG:
		p.parseHashEndFrag()
	case ENDGEO:
		p.parseHashEndGeo()
	default:
		p.write(literal)
		p.write(p.literal)
	}
}

func (p *Parser) parseHashInclude() {
	p.next()
	if p.token == WHITESPACE {
		p.next()
	}

	if p.token != STRING {
		p.unexpectedError()
		return
	}

	includePath := filepath.Join(filepath.Dir(p.path), p.literal)
	src, err := os.Open(includePath)
	if err != nil {
		fmt.Println("Error opening include path: ", includePath, " : ", err)
		return
	}
	defer src.Close()

	fragBuf, vertBuf, geoBuf := new(bytes.Buffer), new(bytes.Buffer), new(bytes.Buffer)
	New(src, includePath, fragBuf, vertBuf, geoBuf).Parse()
	p.writeFrag(fragBuf.Bytes())
	p.writeVert(vertBuf.Bytes())
	p.writeGeo(geoBuf.Bytes())
}

func (p *Parser) parseHashVert() {
	if p.frag || p.geo {
		p.error("Cannot have #vert inside another exclusion.")
	} else {
		p.vert = true
	}
}

func (p *Parser) parseHashFrag() {
	if p.vert || p.geo {
		p.error("Cannot have #frag inside another exclusion.")
	} else {
		p.frag = true
	}
}

func (p *Parser) parseHashGeo() {
	if p.vert || p.frag {
		p.error("Cannot have #geo inside another exclusion.")
	} else {
		p.geo = true
	}
}

func (p *Parser) parseHashEndVert() {
	if !p.vert {
		p.error("Unexpected #endvert")
	}
	p.vert = false
}

func (p *Parser) parseHashEndFrag() {
	if !p.frag {
		p.error("Unexpected #endfrag")
	}
	p.frag = false
}

func (p *Parser) parseHashEndGeo() {
	if !p.geo {
		p.error("Unexpected #endgeo")
	}
	p.geo = false
}

func (p *Parser) parseString() {
	p.write(`"`)
	p.write(p.literal)
	p.write(`"`)
}

func (p *Parser) next() {
	p.token, p.literal = p.scan()
}

func (p *Parser) unexpectedError() {
	p.error(fmt.Sprint("Unexpected '", p.literal))
}

func (p *Parser) error(msg string) {
	fmt.Println("Error: ", msg)
}

func (p *Parser) write(str string) {
	data := []byte(str)
	if p.frag {
		p.writeFrag(data)
	} else if p.vert {
		p.writeVert(data)
	} else if p.geo {
		p.writeGeo(data)
	} else {
		p.writeAll(data)
	}
}

func (p *Parser) writeAll(data []byte) {
	p.writeFrag(data)
	p.writeVert(data)
	p.writeGeo(data)
}

func (p *Parser) writeFrag(data []byte) {
	if p.fragOut != nil {
		p.fragOut.Write(data)
	}
}

func (p *Parser) writeVert(data []byte) {
	if p.vertOut != nil {
		p.vertOut.Write(data)
	}
}

func (p *Parser) writeGeo(data []byte) {
	if p.geoOut != nil {
		p.geoOut.Write(data)
	}
}
