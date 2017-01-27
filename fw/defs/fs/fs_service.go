// Code generated by clubbygen.
// GENERATED FILE DO NOT EDIT
// +build !clubby_strict

package fs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"cesanta.com/common/go/mgrpc"
	"cesanta.com/common/go/mgrpc/frame"
	"cesanta.com/common/go/ourjson"
	"cesanta.com/common/go/ourtrace"
	"github.com/cesanta/errors"
	"golang.org/x/net/trace"
)

var _ = bytes.MinRead
var _ = fmt.Errorf
var emptyMessage = ourjson.RawMessage{}
var _ = ourtrace.New
var _ = trace.New

const ServiceID = "http://mongoose-iot.com/fwFS"

type GetArgs struct {
	Filename *string `json:"filename,omitempty"`
	Len      *int64  `json:"len,omitempty"`
	Offset   *int64  `json:"offset,omitempty"`
}

type GetResult struct {
	Data *string `json:"data,omitempty"`
	Left *int64  `json:"left,omitempty"`
}

type PutArgs struct {
	Append   *bool   `json:"append,omitempty"`
	Data     *string `json:"data,omitempty"`
	Filename *string `json:"filename,omitempty"`
}

type Service interface {
	Get(ctx context.Context, args *GetArgs) (*GetResult, error)
	List(ctx context.Context) ([]string, error)
	Put(ctx context.Context, args *PutArgs) error
}

type Instance interface {
	Call(context.Context, string, *frame.Command) (*frame.Response, error)
}

func NewClient(i Instance, addr string) Service {
	return &_Client{i: i, addr: addr}
}

type _Client struct {
	i    Instance
	addr string
}

func (c *_Client) Get(ctx context.Context, args *GetArgs) (res *GetResult, err error) {
	cmd := &frame.Command{
		Cmd: "FS.Get",
	}

	cmd.Args = ourjson.DelayMarshaling(args)
	if args.Filename == nil {
		return nil, errors.Errorf("Filename is required")
	}
	resp, err := c.i.Call(ctx, c.addr, cmd)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if resp.Status != 0 {
		return nil, errors.Trace(&mgrpc.ErrorResponse{Status: resp.Status, Msg: resp.StatusMsg})
	}

	var r *GetResult
	err = resp.Response.UnmarshalInto(&r)
	if err != nil {
		return nil, errors.Annotatef(err, "unmarshaling response")
	}
	return r, nil
}

func (c *_Client) List(ctx context.Context) (res []string, err error) {
	cmd := &frame.Command{
		Cmd: "FS.List",
	}
	resp, err := c.i.Call(ctx, c.addr, cmd)
	if err != nil {
		return nil, errors.Trace(err)
	}
	if resp.Status != 0 {
		return nil, errors.Trace(&mgrpc.ErrorResponse{Status: resp.Status, Msg: resp.StatusMsg})
	}

	var r []string
	err = resp.Response.UnmarshalInto(&r)
	if err != nil {
		return nil, errors.Annotatef(err, "unmarshaling response")
	}
	return r, nil
}

func (c *_Client) Put(ctx context.Context, args *PutArgs) (err error) {
	cmd := &frame.Command{
		Cmd: "FS.Put",
	}

	cmd.Args = ourjson.DelayMarshaling(args)
	if args.Filename == nil {
		return errors.Errorf("Filename is required")
	}
	resp, err := c.i.Call(ctx, c.addr, cmd)
	if err != nil {
		return errors.Trace(err)
	}
	if resp.Status != 0 {
		return errors.Trace(&mgrpc.ErrorResponse{Status: resp.Status, Msg: resp.StatusMsg})
	}
	return nil
}

//func RegisterService(i *clubby.Instance, impl Service) error {
//s := &_Server{impl}
//i.RegisterCommandHandler("FS.Get", s.Get)
//i.RegisterCommandHandler("FS.List", s.List)
//i.RegisterCommandHandler("FS.Put", s.Put)
//i.RegisterService(ServiceID, _ServiceDefinition)
//return nil
//}

type _Server struct {
	impl Service
}

func (s *_Server) Get(ctx context.Context, src string, cmd *frame.Command) (interface{}, error) {
	var args GetArgs
	if len(cmd.Args) > 0 {
		if err := cmd.Args.UnmarshalInto(&args); err != nil {
			return nil, errors.Annotatef(err, "unmarshaling args")
		}
	}
	if args.Filename == nil {
		return nil, errors.Errorf("Filename is required")
	}
	return s.impl.Get(ctx, &args)
}

func (s *_Server) List(ctx context.Context, src string, cmd *frame.Command) (interface{}, error) {
	return s.impl.List(ctx)
}

func (s *_Server) Put(ctx context.Context, src string, cmd *frame.Command) (interface{}, error) {
	var args PutArgs
	if len(cmd.Args) > 0 {
		if err := cmd.Args.UnmarshalInto(&args); err != nil {
			return nil, errors.Annotatef(err, "unmarshaling args")
		}
	}
	if args.Filename == nil {
		return nil, errors.Errorf("Filename is required")
	}
	return nil, s.impl.Put(ctx, &args)
}

var _ServiceDefinition = json.RawMessage([]byte(`{
  "methods": {
    "Get": {
      "args": {
        "filename": {
          "doc": "Name of the file to read.",
          "type": "string"
        },
        "len": {
          "doc": "Length of chunk to read. If omitted, all available data until the EOF\nwill be read. If (offset + len) is larger than the file size, no\nerror will be returned, and only available data until the EOF will be\nread.\n",
          "type": "integer"
        },
        "offset": {
          "doc": "Offset from the beginning of the file to start reading from.\nIf omitted, 0 is assumed. If the given offset is larger than the file\nsize, no error is returned, and the returned data will be null.\n",
          "type": "integer"
        }
      },
      "doc": "Read a file or a part of file from the device's filesystem.",
      "required_args": [
        "filename"
      ],
      "result": {
        "properties": {
          "data": {
            "doc": "Base64-encoded chunk of data read from the file.",
            "type": "string"
          },
          "left": {
            "doc": "Number of bytes left past the read chunk of data.",
            "type": "integer"
          }
        },
        "type": "object"
      }
    },
    "List": {
      "doc": "List files at the device's filesystem.",
      "result": {
        "items": {
          "doc": "Filename",
          "type": "string"
        },
        "type": "array"
      }
    },
    "Put": {
      "args": {
        "append": {
          "doc": "If true, and if the file with the given filename already exists, the\ndata will be appended to it. Otherwise, the file will be overwritten\nor created.\n",
          "type": "boolean"
        },
        "data": {
          "doc": "Base-64 encoded data to write or append.",
          "type": "string"
        },
        "filename": {
          "doc": "Name of the file to write or append to.",
          "type": "string"
        }
      },
      "doc": "Write or append data to file.",
      "required_args": [
        "filename"
      ]
    }
  },
  "name": "FS",
  "namespace": "http://mongoose-iot.com/fw"
}`))