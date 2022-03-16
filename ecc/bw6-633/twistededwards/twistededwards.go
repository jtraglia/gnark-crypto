/*
Copyright © 2020 ConsenSys

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package twistededwards

import (
	"math/big"

	"github.com/consensys/gnark-crypto/ecc/bw6-633/fr"
)

// CurveParams curve parameters: ax^2 + y^2 = 1 + d*x^2*y^2
type CurveParams struct {
	A, D     fr.Element // in Montgomery form
	Cofactor fr.Element // not in Montgomery form
	Order    big.Int
	Base     PointAffine
}

var edwards CurveParams

// GetEdwardsCurve returns the twisted Edwards curve on BW6-633's Fr
func GetEdwardsCurve() CurveParams {
	// copy to keep Order private
	var res CurveParams

	res.A.Set(&edwards.A)
	res.D.Set(&edwards.D)
	res.Cofactor.Set(&edwards.Cofactor)
	res.Order.Set(&edwards.Order)
	res.Base.Set(&edwards.Base)

	return res
}

func init() {

	edwards.A.SetOne().Neg(&edwards.A)
	edwards.D.SetString("37248940285811842784899494310834635440994424264352085037441815381151934266434102922992043546621")
	edwards.Cofactor.SetUint64(8).ToMont()
	edwards.Order.SetString("4963142838689179791878211236301121218116687802119716497817028544854034649070444389864454748079", 10)

	edwards.Base.X.SetString("37635937024655419978837220647164498012335808680404874556501960268316961933409049243153117555100")
	edwards.Base.Y.SetString("23823085625708063001015413934245381846960101450148849601038571303382730455875805408244170280142")
}

// mulByA multiplies fr.Element by edwards.A
func mulByA(x *fr.Element) {
	x.Neg(x)
}
