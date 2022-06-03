package method

import (
	"fmt"
)

// 伽罗瓦域计算
// 加
func Galois_ADD(x, y int) int {
	return x ^ y
}

// 减
func Galois_SUB(x, y int) int {
	return x ^ y
}

// 乘，不使用速查表
func Galois_MUL_NO_TABLE(x, y, prim int) int {
	res := 0
	for y > 0 {
		if y&1 > 0 {
			res ^= x
		}
		y >>= 1
		x <<= 1
		if prim > 0 && (x&256 > 0) {
			x ^= prim
		}
	}
	return res
}

// 乘，使用速查表
func Galois_MUL(x, y int) int {
	if x == 0 || y == 0 {
		return 0
	}
	return Galois_GF_EXP_256[(Galois_GF_LOG_256[x]+Galois_GF_LOG_256[y])%255]
}

// 除法，使用速查表
func Galois_DIV(x, y int) int {
	if y == 0 {
		// TODO: throw error
		return 0
	}
	if x == 0 {
		return 0
	}
	return Galois_GF_EXP_256[(Galois_GF_LOG_256[x]+255-Galois_GF_LOG_256[y])%255]
}

// 幂，使用速查表
func Galois_POW(x, power int) int {
	return Galois_GF_EXP_256[(Galois_GF_LOG_256[x]*power)%255]
}

// 逆，使用速查表，求 1/x = Galois_DIV(1,x)
func Galois_INVERSE(x int) int {
	return Galois_GF_EXP_256[255-Galois_GF_LOG_256[x]]
}

// 多项式，乘以标量，返回新的多项式
func Galois_POLY_SCALE(p []int, x int) []int {
	result := make([]int, len(p))
	for index, pItem := range p {
		result[index] = Galois_MUL(pItem, x)
	}
	return result
}

// 多项式，加 多项式，返回新的多项式
func Galois_POLY_ADD(p []int, q []int) []int {
	length := len(p)
	if len(q) > length {
		length = len(q)
	}
	result := make([]int, length)
	for index, pItem := range p {
		result[index+length-len(p)] = pItem
	}
	for index, qItem := range q {
		result[index+length-len(q)] ^= qItem
	}
	return result
}

// 多项式，乘 多项式，返回新的多项式
func Galois_POLY_MUL(p []int, q []int) []int {
	length := len(p) + len(q) - 1
	result := make([]int, length)

	for j := range q {
		for i := range p {
			result[i+j] ^= Galois_MUL(p[i], q[j])
		}
	}

	return result
}

// 多项式，计算多项式在特定x的值
func Galois_POLY_EVAL(p []int, x int) int {
	y := p[0]
	for i, item := range p {
		if i == 0 {
			continue
		}
		y = Galois_MUL(y, x) ^ item
	}
	return y
}

func Galois_Table() {
	x := 1
	gf_exp := make([]int, 256)
	gf_log := make([]int, 256)
	for i := 0; i < 256; i++ {
		gf_exp[i] = x
		gf_log[x] = i
		x = Galois_MUL_NO_TABLE(x, 2, 0x11d)
	}
	fmt.Println("gf_exp", gf_exp)
	fmt.Println("gf_log", gf_log)
}

// 都是通过上述方法执行生成的
var Galois_GF_EXP_256 = [512]int{
	1, 2, 4, 8, 16, 32, 64, 128, 29, 58, 116, 232, 205, 135, 19, 38, 76, 152, 45, 90, 180, 117, 234, 201, 143, 3, 6, 12, 24, 48, 96, 192, 157, 39, 78, 156, 37, 74, 148, 53, 106, 212, 181, 119, 238, 193, 159, 35, 70, 140, 5, 10, 20, 40, 80, 160, 93, 186, 105, 210, 185, 111, 222, 161, 95, 190, 97, 194, 153, 47, 94, 188, 101, 202, 137, 15, 30, 60, 120, 240, 253, 231, 211, 187, 107, 214, 177, 127, 254, 225, 223, 163, 91, 182, 113, 226, 217, 175, 67, 134, 17, 34, 68, 136, 13, 26, 52, 104, 208, 189, 103, 206, 129, 31, 62, 124, 248, 237, 199, 147, 59, 118, 236, 197, 151, 51, 102, 204, 133, 23, 46, 92, 184, 109, 218, 169, 79, 158, 33, 66, 132, 21, 42, 84, 168, 77, 154, 41, 82, 164, 85, 170, 73, 146, 57, 114, 228, 213, 183, 115, 230, 209, 191, 99, 198, 145, 63, 126, 252, 229, 215, 179, 123, 246, 241, 255, 227, 219, 171, 75, 150, 49, 98, 196, 149, 55, 110, 220, 165, 87, 174, 65, 130, 25, 50, 100, 200, 141, 7, 14, 28, 56, 112, 224, 221, 167, 83, 166, 81, 162, 89, 178, 121, 242, 249, 239, 195, 155, 43, 86, 172, 69, 138, 9, 18, 36, 72, 144, 61, 122, 244, 245, 247, 243, 251, 235, 203, 139, 11, 22, 44, 88, 176, 125, 250, 233, 207, 131, 27, 54, 108, 216, 173, 71, 142, 1,
}
var Galois_GF_LOG_256 = [256]int{
	0, 255, 1, 25, 2, 50, 26, 198, 3, 223, 51, 238, 27, 104, 199, 75, 4, 100, 224, 14, 52, 141, 239, 129, 28, 193, 105, 248, 200, 8, 76, 113, 5, 138, 101, 47, 225, 36, 15, 33, 53, 147, 142, 218, 240, 18, 130, 69, 29, 181, 194, 125, 106, 39, 249, 185, 201, 154, 9, 120, 77, 228, 114, 166, 6, 191, 139, 98, 102, 221, 48, 253, 226, 152, 37, 179, 16, 145, 34, 136, 54, 208, 148, 206, 143, 150, 219, 189, 241, 210, 19, 92, 131, 56, 70, 64, 30, 66, 182, 163, 195, 72, 126, 110, 107, 58, 40, 84, 250, 133, 186, 61, 202, 94, 155, 159, 10, 21, 121, 43, 78, 212, 229, 172, 115, 243, 167, 87, 7, 112, 192, 247, 140, 128, 99, 13, 103, 74, 222, 237, 49, 197, 254, 24, 227, 165, 153, 119, 38, 184, 180, 124, 17, 68, 146, 217, 35, 32, 137, 46, 55, 63, 209, 91, 149, 188, 207, 205, 144, 135, 151, 178, 220, 252, 190, 97, 242, 86, 211, 171, 20, 42, 93, 158, 132, 60, 57, 83, 71, 109, 65, 162, 31, 45, 67, 216, 183, 123, 164, 118, 196, 23, 73, 236, 127, 12, 111, 246, 108, 161, 59, 82, 41, 157, 85, 170, 251, 96, 134, 177, 187, 204, 62, 90, 203, 89, 95, 176, 156, 169, 160, 81, 11, 245, 22, 235, 122, 117, 44, 215, 79, 174, 213, 233, 230, 231, 173, 232, 116, 214, 244, 234, 168, 80, 88, 175,
}

// 根据纠错符号数量，生成多项式
func RS_Generator_Poly(errorCorrectionSize int) []int {
	g := []int{1}
	for i := 0; i < errorCorrectionSize; i++ {
		g = Galois_POLY_MUL(g, []int{1, Galois_POW(2, i)})
	}
	return g
}

// 证候多项式的计算,如果被扫描消息未损坏，结果应为零
func RS_Calc_Syndromes(msg []int, errCorrectSize int) []int {
	synd := make([]int, errCorrectSize)
	for i := 0; i < errCorrectSize; i++ {
		synd[i] = Galois_POLY_EVAL(msg, Galois_POW(2, i))
	}
	return synd
}

func RS_Encode(msg []int, errCorrectSize int) []int {
	poly := RS_Generator_Poly(errCorrectSize)
	res := make([]int, len(msg)+len(poly)-1)

	for i, item := range msg {
		res[i] = item
	}
	for i := range msg {
		coef := res[i]
		if coef != 0 {
			for j := 1; j < len(poly); j++ {
				res[i+j] ^= Galois_MUL(poly[j], coef)
			}
		}
	}
	for i, item := range msg {
		res[i] = item
	}
	return res
}

// 计算Forney症状
func RS_Forney_Syndromes(synd []int, erase_pos []int, msgSize int) []int {
	erase_pos_reversed := make([]int, len(erase_pos))
	for i, erase := range erase_pos {
		erase_pos_reversed[i] = msgSize - 1 - erase
	}
	fsynd := synd[1:]
	for _, erase := range erase_pos_reversed {
		x := Galois_POW(2, erase)
		for j := 0; j < len(fsynd)-1; j++ {
			fsynd[j] = Galois_MUL(fsynd[j], x) ^ fsynd[j+1]
		}
	}
	return fsynd
}

// 使用 Berlekamp-Massey 计算错误定位器多项式
func RS_Find_Error_Locator(synd []int, nsym int, erase_count int) []int {
	// err_loc := []int{1}
	// auxi_poly := []int{1}
	// L := 0
	// for r := 0; r < nsym; r++ {
	// 	delta := synd[r]
	// 	for j := 1; j < len(err_loc); j++ {
	// 		delta = Galois_SUB(delta, Galois_MUL(err_loc[len(err_loc)-(j+1)], synd[r-j]))
	// 	}

	// 	auxi_poly = append(auxi_poly, 0)
	// 	if delta != 0 {
	// 		old_err_loc := err_loc
	// 		err_loc = Galois_POLY_ADD(err_loc, Galois_POLY_SCALE(auxi_poly, delta))
	// 		if 2*L <= r-1 {
	// 			L = r - L
	// 			auxi_poly = Galois_POLY_SCALE(old_err_loc, Galois_INVERSE(delta))
	// 		}
	// 	}
	// }
	// return err_loc
	err_loc := []int{1}
	old_loc := []int{1}
	synd_shift := len(synd) - nsym
	// if len(synd) > nsym {
	// 	synd_shift = len(synd) - nsym
	// }
	for i := 0; i < nsym-erase_count; i++ {
		// fmt.Println("xxi", i, nsym, synd, synd_shift)
		K := i + synd_shift
		var delta int
		if K < 0 {
			delta = synd[len(synd)+K]
		} else {
			delta = synd[K]
		}

		// fmt.Println("delta", i, delta, K, synd)
		for j := 1; j < len(err_loc); j++ {
			syncIndex := K - j
			if syncIndex < 0 {
				syncIndex = len(synd) + syncIndex
			}
			delta ^= Galois_MUL(err_loc[len(err_loc)-(j+1)], synd[syncIndex])
		}
		old_loc = append(old_loc, 0)
		if delta != 0 {
			if len(old_loc) > len(err_loc) {
				new_loc := Galois_POLY_SCALE(old_loc, delta)
				old_loc = Galois_POLY_SCALE(err_loc, Galois_INVERSE(delta))
				err_loc = new_loc
			}
			err_loc = Galois_POLY_ADD(err_loc, Galois_POLY_SCALE(old_loc, delta))
		}
	}
	for len(err_loc) > 0 && err_loc[0] == 0 {
		err_loc = err_loc[1:]
	}
	errs := len(err_loc) - 1
	if (errs-erase_count)*2+erase_count > nsym {
		// TODO: Error Too many errors to correct
	}
	return err_loc
}

//
func RS_Find_Errors(err_loc []int, msgLength int) []int {
	fmt.Println("err_loc", err_loc, msgLength)
	err_pos := []int{}
	for i := 0; i < msgLength; i++ {
		fmt.Println("err_loc", err_loc, Galois_POW(2, i))
		if Galois_POLY_EVAL(err_loc, Galois_POW(2, i)) == 0 {
			err_pos = append(err_pos, msgLength-1-i)
		}
	}
	return err_pos
}

// 错误纠正
func RS_Error_Correct(msg []int, nsym int, erase_pos []int) []int {
	fmt.Println("msg", msg, nsym)
	output := msg[:]
	for _, pos := range erase_pos {
		output[pos] = 0
	}
	synd := RS_Calc_Syndromes(output, nsym)
	fmt.Println("synd", synd)
	isAllSyndIsZero := true
	for _, syndNum := range synd {
		if syndNum != 0 {
			isAllSyndIsZero = false
		}
	}
	// 没有错误
	if isAllSyndIsZero {
		return output
	}
	fsynd := RS_Forney_Syndromes(synd, erase_pos, len(output))
	fmt.Println("fsynd", synd, fsynd)
	err_loc := RS_Find_Error_Locator(fsynd, nsym, len(erase_pos))
	err_loc_reverse := make([]int, len(err_loc))
	for index, item := range err_loc {
		err_loc_reverse[len(err_loc)-index-1] = item
	}
	err_pos := RS_Find_Errors(err_loc_reverse, len(output))

	// output = rs_correct_errata(output, synd, (erase_pos + err_pos))
	fmt.Println("err_pos", err_pos)
	return output
}
