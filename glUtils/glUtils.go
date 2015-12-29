// augment Sylvester some
package glUtils

import (
    "math"
    "golang.org/x/image/math/f32"
)

func TranslateMatrix(m f32.Mat4, v f32.Vec3) f32.Mat4 {
    vv := Identity();
    vv[3] = v[0];
    vv[7] = v[1];
    vv[11] = v[2];
    return X4(m, vv);
}

func Flatten(m *f32.Mat4) []float32 {
    a := make([]float32, 16);
    for i := 0; i < 16; i++ {
        r := i % 4;
        c := i / 4;
        a[r * 4 + c] = m[i];
    }
    return a;
}

func Identity() f32.Mat4 {
    return f32.Mat4 { 1,0,0,0, 0,1,0,0, 0,0,1,0, 0,0,0,1 }
}

//
// gluPerspective
//
func MakePerspective(fovy float32, aspect float32, znear float32, zfar float32) f32.Mat4 {
    ymax := znear * float32(math.Tan(float64(fovy) * math.Pi / 360.0));
    ymin := -ymax;
    xmin := ymin * aspect;
    xmax := ymax * aspect;

    return MakeFrustum(xmin, xmax, ymin, ymax, znear, zfar);
}

//
// glFrustum
//
func MakeFrustum(left float32, right float32, bottom float32, top float32, znear float32, zfar float32) f32.Mat4 {
    X := 2*znear/(right-left);
    Y := 2*znear/(top-bottom);
    A := (right+left)/(right-left);
    B := (top+bottom)/(top-bottom);
    C := -(zfar+znear)/(zfar-znear);
    D := -2*zfar*znear/(zfar-znear);

    return f32.Mat4 {
        X, 0, A, 0,
        0, Y, B, 0,
        0, 0, C, D,
        0, 0, -1, 0 };
}

//
// Multiply 2 4x4 matrices, return a new matrix
//
func X4(m1, m2 f32.Mat4) f32.Mat4 {
    i := 4;
    nj := 4;
    j := 0;
    cols := 4
    c := 0
    var sum float32;

    sum = 0.0;

    result := Identity();

    for i != 0 { i--; j = nj;
        for j != 0 { j--; c = cols;
            sum = 0;
            for c != 0 { c--;
              sum += m1[i*4+c] * m2[c*4+j];
            }
            result[i*4+j] = sum;
        }
    }

    return result;
}
