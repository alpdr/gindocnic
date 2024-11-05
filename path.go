package gindocnic

func filterNonAlphaNumeric(s string) string {
	bytes := []byte(s)
	res := make([]byte, 0)
	for _, c := range bytes {
		if isAlphaNumeric(c) {
			res = append(res, c)
		}
	}
	return string(res)
}

// Open API形式のパス /v1/profiles/{id}/*actionsから*actionsの部分を探します。
func findStarParams(openAPIPath string) map[string]bool {
	seq := []rune(openAPIPath)
	n := len(seq)
	res := make(map[string]bool, 0)
	i := 0
	for i < n {
		if seq[i] == '*' {
			i++
			param := make([]rune, 0)
			for i < n && seq[i] != '/' {
				param = append(param, seq[i])
				i++
			}
			res[string(param)] = true
		}
		i++
	}
	return res
}

// ginのパスをOpenAPIのパスに変換
// たとえば、/v1/profiles/:id -> /v1/profiles/{id}に変換します。
func makeGinToOpenAPIPath(ginPath string) string {
	seq := []rune(ginPath)
	n := len(seq)
	res := make([]rune, 0)

	i := 0
	for i < n {
		if seq[i] == ':' {
			res = append(res, '{')
			i++
			for i < n && seq[i] != '/' {
				res = append(res, seq[i])
				i++
			}
			res = append(res, '}')
		} else {
			res = append(res, seq[i])
			i++
		}
	}
	return string(res)
}

func isAlphaNumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
}
