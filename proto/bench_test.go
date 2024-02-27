// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package proto_test

import (
	"flag"
	"fmt"
	"reflect"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/protoadapt"
)

// The results of these microbenchmarks are unlikely to correspond well
// to real world performance. They are mainly useful as a quick check to
// detect unexpected regressions and for profiling specific cases.

var (
	allowPartial = flag.Bool("allow_partial", false, "set AllowPartial")
)

// BenchmarkEncode benchmarks encoding all the test messages.
func BenchmarkEncode(b *testing.B) {
	s := PBSerializationNew{}
	for _, test := range testValidMessages {
		for _, want := range test.decodeTo {
			b.Run(fmt.Sprintf("%s (%T)", test.desc, want), func(b *testing.B) {
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						_, err := s.Marshal(want)
						if err != nil && !test.partial {
							b.Fatal(err)
						}
					}
				})
			})
		}
	}
}

// BenchmarkDecode benchmarks decoding all the test messages.
func BenchmarkDecode(b *testing.B) {
	s := PBSerializationNew{}
	for _, test := range testValidMessages {
		for _, want := range test.decodeTo {
			b.Run(fmt.Sprintf("%s (%T)", test.desc, want), func(b *testing.B) {
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						m := reflect.New(reflect.TypeOf(want).Elem()).Interface().(proto.Message)
						err := s.Unmarshal(test.wire, m)
						if err != nil && !test.partial {
							b.Fatal(err)
						}
					}
				})
			})
		}
	}
}

// // BenchmarkEncode benchmarks encoding all the test messages.
// func BenchmarkEncodeNew(b *testing.B) {
//     s := PBSerializationNew{}
//     for _, test := range testValidMessages {
//         for _, want := range test.decodeTo {
//             b.Run(fmt.Sprintf("%s (%T)", test.desc, want), func(b *testing.B) {
//                 b.RunParallel(func(pb *testing.PB) {
//                     for pb.Next() {
//                         _, err := s.Marshal(want)
//                         if err != nil && !test.partial {
//                             b.Fatal(err)
//                         }
//                     }
//                 })
//             })
//         }
//     }
// }
//
// // BenchmarkDecode benchmarks decoding all the test messages.
// func BenchmarkDecodeNew(b *testing.B) {
//     s := PBSerializationNew{}
//     for _, test := range testValidMessages {
//         for _, want := range test.decodeTo {
//             b.Run(fmt.Sprintf("%s (%T)", test.desc, want), func(b *testing.B) {
//                 b.RunParallel(func(pb *testing.PB) {
//                     for pb.Next() {
//                         m := reflect.New(reflect.TypeOf(want).Elem()).Interface().(proto.Message)
//                         err := s.Unmarshal(test.wire, m)
//                         if err != nil && !test.partial {
//                             b.Fatal(err)
//                         }
//                     }
//                 })
//             })
//         }
//     }
// }

// PBSerializationOld provides protobuf serialization mode.
type PBSerializationOld struct{}

// Unmarshal deserializes the in bytes into body.
func (s *PBSerializationOld) Unmarshal(in []byte, body interface{}) error {
	msg, ok := body.(proto.Message)
	if !ok {
		return fmt.Errorf("failed to unmarshal body: expected proto.Message, got %T", body)
	}
	return proto.Unmarshal(in, msg)
}

// Marshal returns the serialized bytes in protobuf protocol.
func (s *PBSerializationOld) Marshal(body interface{}) ([]byte, error) {
	msg, ok := body.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("failed to marshal body: expected proto.Message, got %T", body)
	}
	return proto.Marshal(msg)
}

// PBSerializationNew provides protobuf serialization mode.
type PBSerializationNew struct{}

// Unmarshal deserializes the in bytes into body.
func (s *PBSerializationNew) Unmarshal(in []byte, body interface{}) error {
	msgV2, ok := body.(protoadapt.MessageV2)
	if !ok {
		msgV1, ok := body.(protoadapt.MessageV1)
		if !ok {
			return fmt.Errorf("failed to unmarshal body: expected github.com/golang/protobuf/proto.Message "+
				"or google.golang.org/protobuf/proto.Message,  got %T", body)
		}
		msgV2 = protoadapt.MessageV2Of(msgV1)
	}
	return proto.Unmarshal(in, msgV2)
}

// Marshal returns the serialized bytes in protobuf protocol.
func (s *PBSerializationNew) Marshal(body interface{}) ([]byte, error) {
	msgV2, ok := body.(protoadapt.MessageV2)
	if !ok {
		msgV1, ok := body.(protoadapt.MessageV1)
		if !ok {
			return nil, fmt.Errorf("failed to marshal body: expected github.com/golang/protobuf/proto.Message "+
				"or google.golang.org/protobuf/proto.Message,  got %T", body)
		}
		msgV2 = protoadapt.MessageV2Of(msgV1)
	}
	return proto.Marshal(msgV2)
}
