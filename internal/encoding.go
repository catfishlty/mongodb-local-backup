package internal

import "golang.org/x/text/encoding/simplifiedchinese"

type Charset string

const (
    UFT8    = Charset("UTF-8")
    GB18030 = Charset("GB18030")
)

func ConvertByteToString(data []byte, charset Charset) string {
    switch charset {
    case GB18030:
        decodeBytes, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(data)
        return string(decodeBytes)
    case UFT8:
        fallthrough
    default:
    }
    return string(data)
}
