package method

import (
	"errors"
	"fmt"
)

// 伽罗瓦域计算
// refs: https://en.wikiversity.org/wiki/Reed%E2%80%93Solomon_codes_for_coders
// refs: https://github.com/paukerspinner/encode-decode_ReedSolomonCode
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
func Galois_MUL(x, y int, table Galois_GF) int {
	if x == 0 || y == 0 {
		return 0
	}
	return table.exp[(table.log[x]+table.log[y])%table.maxIndex]
}

// 除法，使用速查表
func Galois_DIV(x, y int, table Galois_GF) int {
	if y == 0 {
		// TODO: throw error
		return 0
	}
	if x == 0 {
		return 0
	}
	return table.exp[(table.log[x]+255-table.log[y])%table.maxIndex]
}

// 幂，使用速查表
func Galois_POW(x, power int, table Galois_GF) int {
	index := (table.log[x] * power) % table.maxIndex
	// fmt.Println("pow", index, x, power, table.log[x]*power)
	if index < 0 {
		index += table.maxIndex
	}

	return table.exp[index]
}

// 逆，使用速查表，求 1/x = Galois_DIV(1,x)
func Galois_INVERSE(x int, table Galois_GF) int {
	return table.exp[table.maxIndex-table.log[x]]
}

// 多项式，乘以标量，返回新的多项式
func Galois_POLY_SCALE(p []int, x int, table Galois_GF) []int {
	result := make([]int, len(p))
	for index, pItem := range p {
		result[index] = Galois_MUL(pItem, x, table)
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
func Galois_POLY_MUL(p []int, q []int, table Galois_GF) []int {
	length := len(p) + len(q) - 1
	result := make([]int, length)

	for j := range q {
		for i := range p {
			result[i+j] ^= Galois_MUL(p[i], q[j], table)
		}
	}

	return result
}

// 多项式，除法
func Galois_POLY_DIV(dividend, divisor []int, table Galois_GF) ([]int, []int) {
	msg_out := dividend[:]
	length := len(dividend) - (len(divisor) - 1)
	for i := 0; i < length; i++ {
		coef := msg_out[i]
		if coef != 0 {
			for j := 1; j < len(divisor); j++ {
				if divisor[j] != 0 {
					msg_out[i+j] ^= Galois_MUL(divisor[j], coef, table)
				}
			}
		}
	}
	separator := -(len(divisor) - 1)
	if separator >= 0 {
		return msg_out[:separator], msg_out[separator:]
	}
	msg_out_len := len(msg_out)
	return msg_out[:msg_out_len+separator], msg_out[msg_out_len+separator:]
}

// 多项式，计算多项式在特定x的值
func Galois_POLY_EVAL(p []int, x int, table Galois_GF) int {
	y := p[0]
	for i := 1; i < len(p); i++ {
		y = Galois_MUL(y, x, table) ^ p[i]
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
}

type Galois_GF struct {
	maxIndex int
	exp      []int
	log      []int
}

// 都是通过上述方法执行生成的
var Galois_GF_EXP_256 = []int{
	1, 2, 4, 8, 16, 32, 64, 128, 29, 58, 116, 232, 205, 135, 19, 38, 76, 152, 45, 90, 180, 117, 234, 201, 143, 3, 6, 12, 24, 48, 96, 192, 157, 39, 78, 156, 37, 74, 148, 53, 106, 212, 181, 119, 238, 193, 159, 35, 70, 140, 5, 10, 20, 40, 80, 160, 93, 186, 105, 210, 185, 111, 222, 161, 95, 190, 97, 194, 153, 47, 94, 188, 101, 202, 137, 15, 30, 60, 120, 240, 253, 231, 211, 187, 107, 214, 177, 127, 254, 225, 223, 163, 91, 182, 113, 226, 217, 175, 67, 134, 17, 34, 68, 136, 13, 26, 52, 104, 208, 189, 103, 206, 129, 31, 62, 124, 248, 237, 199, 147, 59, 118, 236, 197, 151, 51, 102, 204, 133, 23, 46, 92, 184, 109, 218, 169, 79, 158, 33, 66, 132, 21, 42, 84, 168, 77, 154, 41, 82, 164, 85, 170, 73, 146, 57, 114, 228, 213, 183, 115, 230, 209, 191, 99, 198, 145, 63, 126, 252, 229, 215, 179, 123, 246, 241, 255, 227, 219, 171, 75, 150, 49, 98, 196, 149, 55, 110, 220, 165, 87, 174, 65, 130, 25, 50, 100, 200, 141, 7, 14, 28, 56, 112, 224, 221, 167, 83, 166, 81, 162, 89, 178, 121, 242, 249, 239, 195, 155, 43, 86, 172, 69, 138, 9, 18, 36, 72, 144, 61, 122, 244, 245, 247, 243, 251, 235, 203, 139, 11, 22, 44, 88, 176, 125, 250, 233, 207, 131, 27, 54, 108, 216, 173, 71, 142, 1,
}
var Galois_GF_LOG_256 = []int{
	0, 255, 1, 25, 2, 50, 26, 198, 3, 223, 51, 238, 27, 104, 199, 75, 4, 100, 224, 14, 52, 141, 239, 129, 28, 193, 105, 248, 200, 8, 76, 113, 5, 138, 101, 47, 225, 36, 15, 33, 53, 147, 142, 218, 240, 18, 130, 69, 29, 181, 194, 125, 106, 39, 249, 185, 201, 154, 9, 120, 77, 228, 114, 166, 6, 191, 139, 98, 102, 221, 48, 253, 226, 152, 37, 179, 16, 145, 34, 136, 54, 208, 148, 206, 143, 150, 219, 189, 241, 210, 19, 92, 131, 56, 70, 64, 30, 66, 182, 163, 195, 72, 126, 110, 107, 58, 40, 84, 250, 133, 186, 61, 202, 94, 155, 159, 10, 21, 121, 43, 78, 212, 229, 172, 115, 243, 167, 87, 7, 112, 192, 247, 140, 128, 99, 13, 103, 74, 222, 237, 49, 197, 254, 24, 227, 165, 153, 119, 38, 184, 180, 124, 17, 68, 146, 217, 35, 32, 137, 46, 55, 63, 209, 91, 149, 188, 207, 205, 144, 135, 151, 178, 220, 252, 190, 97, 242, 86, 211, 171, 20, 42, 93, 158, 132, 60, 57, 83, 71, 109, 65, 162, 31, 45, 67, 216, 183, 123, 164, 118, 196, 23, 73, 236, 127, 12, 111, 246, 108, 161, 59, 82, 41, 157, 85, 170, 251, 96, 134, 177, 187, 204, 62, 90, 203, 89, 95, 176, 156, 169, 160, 81, 11, 245, 22, 235, 122, 117, 44, 215, 79, 174, 213, 233, 230, 231, 173, 232, 116, 214, 244, 234, 168, 80, 88, 175,
}
var Galois_GF_256 = Galois_GF{
	maxIndex: 255,
	exp:      Galois_GF_EXP_256,
	log:      Galois_GF_LOG_256,
}

var Galois_GF_EXP_16 = []int{
	1, 2, 4, 8, 3, 6, 12, 11, 5, 10, 7, 14, 15, 13, 9, 1,
}

var Galois_GF_LOG_16 = []int{
	0, 15, 1, 4, 2, 8, 5, 10, 3, 14, 9, 7, 6, 13, 11, 12,
}

var Galois_GF_15 = Galois_GF{
	maxIndex: 15,
	exp:      Galois_GF_EXP_16,
	log:      Galois_GF_LOG_16,
}

// 根据纠错符号数量，生成多项式
func RS_Generator_Poly(errorCorrectionSize int, table Galois_GF) []int {
	g := []int{1}
	for i := 0; i < errorCorrectionSize; i++ {
		g = Galois_POLY_MUL(g, []int{1, Galois_POW(2, i, table)}, table)
	}
	return g
}

// 证候多项式的计算,如果被扫描消息未损坏，结果应为零
func RS_Calc_Syndromes(msg []int, errCorrectSize int, table Galois_GF) []int {
	synd := make([]int, errCorrectSize)
	for i := 0; i < errCorrectSize; i++ {
		synd[i] = Galois_POLY_EVAL(msg, Galois_POW(2, i, table), table)
	}
	return ConcatArray([]int{0}, synd)
}

func RS_Encode(msg []int, errCorrectSize int, table Galois_GF) []int {
	poly := RS_Generator_Poly(errCorrectSize, table)
	res := make([]int, len(msg)+len(poly)-1)

	for i, item := range msg {
		res[i] = item
	}
	for i := range msg {
		coef := res[i]
		if coef != 0 {
			for j := 1; j < len(poly); j++ {
				res[i+j] ^= Galois_MUL(poly[j], coef, table)
			}
		}
	}
	for i, item := range msg {
		res[i] = item
	}
	return res
}

// 计算Forney症状
func RS_Forney_Syndromes(synd []int, erase_pos []int, msgSize int, table Galois_GF) []int {
	erase_pos_reversed := make([]int, len(erase_pos))
	for i, erase := range erase_pos {
		erase_pos_reversed[i] = msgSize - 1 - erase
	}
	fsynd := synd[1:]
	for _, erase := range erase_pos_reversed {
		x := Galois_POW(2, erase, table)
		for j := 0; j < len(fsynd)-1; j++ {
			fsynd[j] = Galois_MUL(fsynd[j], x, table) ^ fsynd[j+1]
		}
	}
	return fsynd
}

// 使用 Berlekamp-Massey 计算错误定位器多项式
func RS_Find_Error_Locator(synd []int, nsym int, erase_count int, table Galois_GF) []int {
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
		K := i + synd_shift
		var delta int
		if K < 0 {
			delta = synd[len(synd)+K]
		} else {
			delta = synd[K]
		}

		for j := 1; j < len(err_loc); j++ {
			syncIndex := K - j
			if syncIndex < 0 {
				syncIndex = len(synd) + syncIndex
			}
			delta ^= Galois_MUL(err_loc[len(err_loc)-(j+1)], synd[syncIndex], table)
		}
		old_loc = append(old_loc, 0)
		if delta != 0 {
			if len(old_loc) > len(err_loc) {
				new_loc := Galois_POLY_SCALE(old_loc, delta, table)
				old_loc = Galois_POLY_SCALE(err_loc, Galois_INVERSE(delta, table), table)
				err_loc = new_loc
			}
			err_loc = Galois_POLY_ADD(err_loc, Galois_POLY_SCALE(old_loc, delta, table))
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
func RS_Find_Errors(err_loc []int, msgLength int, table Galois_GF) []int {
	err_pos := []int{}
	for i := 0; i < msgLength; i++ {
		if Galois_POLY_EVAL(err_loc, Galois_POW(2, i, table), table) == 0 {
			err_pos = append(err_pos, msgLength-1-i)
		}
	}
	return err_pos
}

// 错误纠正
func RS_Error_Correct(msg []int, nsym int, erase_pos []int, table Galois_GF) ([]int, error) {
	var err error
	output := msg[:]
	for _, pos := range erase_pos {
		output[pos] = 0
	}
	synd := RS_Calc_Syndromes(output, nsym, table)
	fmt.Println("synd", msg, synd)
	// 没有错误
	if CheckAllIsNum(synd, 0) {
		return output, nil
	}
	fsynd := RS_Forney_Syndromes(synd, erase_pos, len(output), table)
	err_loc := RS_Find_Error_Locator(fsynd, nsym, len(erase_pos), table)
	err_loc_reverse := ReverseArray(err_loc)
	err_pos := RS_Find_Errors(err_loc_reverse, len(output), table)
	if len(err_pos) == 0 {
		return msg, errors.New("Could not locate error")
	}
	allErrPos := ConcatArray(erase_pos, err_pos)
	fmt.Println("output", synd, output, err_pos)
	output, err = RS_Correct_Errata(output, synd, allErrPos, table)
	fmt.Println("output2", output)
	if err != nil {
		return msg, err
	}
	synd = RS_Calc_Syndromes(output, nsym, table)
	if !CheckAllIsNum(synd, 0) {
		fmt.Println("synd", synd, output)
		return msg, errors.New("Could not correct message")
	}
	return output, nil
}

func RS_Correct_Errata(msg []int, synd []int, err_pos []int, table Galois_GF) ([]int, error) {
	coef_pos := make([]int, len(err_pos))
	for index, p := range err_pos {
		coef_pos[index] = len(msg) - 1 - p
	}
	err_loc := RS_Find_Errata_Locator(coef_pos, table)
	err_eval := ReverseArray(RS_Find_Error_Evaluator(ReverseArray(synd), err_loc, len(err_loc)-1, table))
	X := []int{}
	for i := 0; i < len(coef_pos); i++ {
		l := 255 - coef_pos[i]
		X = append(X, Galois_POW(2, table.maxIndex-l, table))
	}
	E := make([]int, len(msg))
	Xlength := len(X)
	for i, Xi := range X {
		Xi_inv := Galois_INVERSE(Xi, table)
		err_loc_prime_tmp := []int{}
		for j := 0; j < Xlength; j++ {
			if j != i {
				err_loc_prime_tmp = append(err_loc_prime_tmp, Galois_SUB(1, Galois_MUL(Xi_inv, X[j], table)))
			}
		}
		err_loc_prime := 1
		for _, coef := range err_loc_prime_tmp {
			err_loc_prime = Galois_MUL(err_loc_prime, coef, table)
		}
		y := Galois_POLY_EVAL(ReverseArray(err_eval), Xi_inv, table)
		y = Galois_MUL(Galois_POW(Xi, 1, table), y, table)
		magnitude := Galois_DIV(y, err_loc_prime, table)
		E[err_pos[i]] = magnitude
	}

	msg = Galois_POLY_ADD(msg, E)
	return msg, nil
}

func RS_Find_Error_Evaluator(synd []int, err_loc []int, nsym int, table Galois_GF) []int {
	nsymArr := make([]int, nsym+1)
	_, remainder := Galois_POLY_DIV(Galois_POLY_MUL(synd, err_loc, table), ConcatArray([]int{0}, nsymArr), table)
	return remainder
}

func RS_Find_Errata_Locator(e_pos []int, table Galois_GF) []int {
	e_loc := []int{1}
	for _, i := range e_pos {
		e_loc = Galois_POLY_MUL(e_loc, Galois_POLY_ADD([]int{1}, []int{Galois_POW(2, i, table), 0}), table)
	}
	return e_loc
}

func BCH_Check_Format(fmt int) int {
	g := 0x537
	for i := 4; i > -1; i-- {
		if fmt&(1<<(i+10)) > 0 {
			fmt ^= g << i
		}
	}
	return fmt
}

func BCH_Hamming_Weight(x int) int {
	weight := 0
	for x > 0 {
		weight += x & 1
		x >>= 1
	}
	return weight
}

func BCH_Decode_Format(format int) (int, int) {
	best_fmt := -1
	bestCode := format
	best_dist := 15
	for test_fmt := 0; test_fmt < 32; test_fmt++ {
		test_code := (test_fmt << 10) ^ BCH_Check_Format(test_fmt<<10)
		test_dist := BCH_Hamming_Weight(format ^ test_code)
		if test_dist < best_dist {
			bestCode =  test_code
			best_dist = test_dist
			best_fmt = test_fmt
		} else if test_dist == best_dist {
			best_fmt = -1
		}
	}
	return best_fmt, bestCode
}
