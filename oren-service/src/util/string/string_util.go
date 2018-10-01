package string

import (
	constant "github.com/fadhilfcr/oren-service/src/util"
	"fmt"
	"strconv"
	"errors"
	"strings"
)

func TableIdFormatter(tag string,no int)(string,error){
	var lenId int
	var err error
	strNo := strconv.Itoa(no)
	id := ""

	switch tag {
	case strings.Trim(constant.TAG_USER," ") :
		lenId = 4
	case strings.Trim(constant.TAG_ADDRESS," ") :
		lenId = 2
	case strings.Trim(constant.TAG_BUSANA," ") :
		lenId = 4
	case strings.Trim(constant.TAG_TRANSAKSI," ") :
		lenId = 4
	case strings.Trim(constant.TAG_REKENING," ") :
		lenId = 4
	case strings.Trim(constant.TAG_RATING," ") :
		lenId = 4
	default:
		lenId = 1
	}

	if len(strNo) <= lenId {
		lenId = lenId - len(strNo)
		for i := 0; i < lenId; i++ {
			id = fmt.Sprintf("%s%s",id,"0")
		}

		id = fmt.Sprintf("%s%s%d",tag,id,no)
	}else {
		err = errors.New("No : out of index")
	}

	return id,err
}
