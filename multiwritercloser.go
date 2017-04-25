// Most of this code is copied from the golang source code
// https://golang.org/src/io/multi.go?m=text
// Copyright 2010 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package multiwritercloser

import "io"

type multiWriterCloser struct {
    writers []io.Writer
}

func (t *multiWriterCloser) Write(p []byte) (n int, err error) {
    for _, w := range t.writers {
        n, err = w.Write(p)
        if err != nil {
            return
        }
        if n != len(p) {
            err = io.ErrShortWrite
            return
        }
    }
    return len(p), nil
}

// best efforts close calls
func (t *multiWriterCloser) Close() error {
    for _, w := range t.writers {
        if closer, ok := w.(io.WriteCloser); ok {
            if err := closer.Close(); err != nil {
                return err
            }
        }
    }
    return nil
}

// MultiWriterCloser creates a writecloser that duplicates its writes to all the
// provided writers, similar to the Unix tee(1) command.
func MultiWriterCloser(writers ...io.Writer) io.WriteCloser {
    w := make([]io.Writer, len(writers))
    copy(w, writers)
    return &multiWriterCloser{w}
}
