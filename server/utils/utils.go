package utils

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/constant"
	"bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/pb"
	authPB "bitbucket.bri.co.id/scm/bricams-addons/qcash-template-service/server/service/stubs/auth"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ArrayToMap(val []string) map[string]string {
	companyOwner := map[string]string{}

	if len(val) > 0 {
		for _, value := range val {
			companyOwner[value] = "1"
		}
	}

	return companyOwner
}

func MapToArrayString(val map[string]string) []string {
	companyOwner := []string{}

	if len(val) > 0 {
		for key, value := range val {
			if value == "true" || value == "t" || value == "1" {
				companyOwner = append(companyOwner, key)
			}
		}

	}

	return companyOwner
}

func MapToJSONString(val map[string]string) string {
	jsonString := "{}"

	if len(val) > 0 {
		b := new(bytes.Buffer)
		for key, value := range val {
			if b.Len() == 0 {
				fmt.Fprintf(b, "\"%s\":\"%s\"", key, value)
			} else {
				fmt.Fprintf(b, ", \"%s\":\"%s\"", key, value)
			}
		}

		jsonString = fmt.Sprintf("{ %s }", b.String())
	}

	return jsonString
}

func JsonStringToMapString(val string) (map[string]string, error) {
	x := map[string]string{}

	err := json.Unmarshal([]byte(val), &x)
	if err != nil {
		return x, err
	}

	return x, err

}

func ContextError(ctx context.Context) error {
	switch ctx.Err() {
	case context.Canceled:
		return status.Error(codes.Canceled, "request is canceled")
	case context.DeadlineExceeded:
		return status.Error(codes.DeadlineExceeded, "deadline is exceeded")
	default:
		return nil
	}
}

// validate string length
func ValidateStringLengthEqual(str string, length int) bool {
	return len(str) == length
}

// validate string length
func ValidateStringLengthLessThan(str string, length int) bool {
	return len(str) < length
}

// validate string length
func ValidateStringLengthRange(str string, min, max int) bool {
	return len(str) >= min && len(str) <= max
}

// email format validation
func ValidateStringEmail(email string) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$`).MatchString(email)
}

// validate country code
func ValidateCountryCode(countryCode string) bool {
	return regexp.MustCompile(`^[A-Z]{2}$`).MatchString(countryCode)
}

// alphanumeric validation
func ValidateStringAlphanumeric(str string) bool {
	return regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString(str)
}

// phone number validation
func ValidatePhoneNumberNumberOnly(str string) bool {
	return regexp.MustCompile(`^[0-9]+$`).MatchString(str)
}

// phone number validation
func ValidatePhoneNumber(str string) bool {
	return regexp.MustCompile(`^[0-9]+$`).MatchString(str)
}

func ValidateISOCurrency(currency string) bool {
	return regexp.MustCompile(`^[A-Z]{3}$`).MatchString(currency)
}

func CompareValueDates(date1, date2 *timestamppb.Timestamp) bool {
	// if diffrent day return fales
	if date1 == nil || date2 == nil {
		return false
	}
	d1y, d1m, d1d := date1.AsTime().Date()
	d2y, d2m, d2d := date2.AsTime().Date()

	return d1y == d2y && d1m == d2m && d1d == d2d
}

func GenerateRandomString(length int) (string, error) {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	randomString := ""
	alphabet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	alphabetLength := len(alphabet)
	for _, b := range randomBytes {
		randomString += string(alphabet[int(b)%alphabetLength])
	}
	return randomString, nil
}

func ContainsSubs(s []string, e string) bool {
	for _, a := range s {
		c := strings.Contains(e, a)
		if c {
			return true
		}
	}
	return false
}

func RemoveCharsRegex(str string) string {
	// Compile the regular expression pattern
	re := regexp.MustCompile(`[.,_&'\\/:@-]`)

	// Remove all matches of the pattern from the string
	str = re.ReplaceAllString(str, " ")

	// Trim any leading or trailing spaces
	str = strings.Join(strings.Fields(str), " ")
	return str
}

func CreateNewCTX(ctx context.Context) context.Context {

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		return metadata.NewOutgoingContext(context.Background(), md)
	}

	return nil
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func UnWrapError(err error) string {
	// rpc error: code = Canceled desc = context canceled
	if strings.Contains(fmt.Sprintf("%v", err), "rpc error:") {
		values := strings.Split(fmt.Sprintf("%v", err), " = ")
		return values[len(values)-1]
	}
	return fmt.Sprintf("%v", err)
}

func ReplaceIfEmpty(str string, fallback string) string {
	if str == "" {
		return fallback
	} else {
		return str
	}
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func CalculateTotalPage(limit uint64, totalRows uint64) uint64 {
	var totalPage uint64 = 1
	if limit > 0 {
		rowMod := totalRows % uint64(limit)
		if rowMod != 0 {
			rowMod = 1
		}
		totalPage = (totalRows / uint64(limit)) + rowMod
	}

	return totalPage
}

func ReplaceStrToDate(date string) *time.Time {

	dateReplace, errReplace := time.Parse("02-01-2006", date)
	if errReplace != nil {
		log.Printf("Error parsing date : %v", errReplace)
		return nil
	}

	return &dateReplace
}

func ReplaceStrToInt(status string) *int {

	statusResult, errReplace := strconv.Atoi(status)
	if errReplace != nil {
		log.Printf("Error parsing status : %v", errReplace)
		return nil
	}

	return &statusResult
}

func GenerateBatchReffNo(prefix string) string {
	return fmt.Sprintf("%s%d", prefix, time.Now().UnixMilli())
}

func PaddedDateString(value string, dateDelimiter string) string {
	date := strings.Split(value, dateDelimiter)
	dayInt, _ := strconv.Atoi(date[0])
	day := PaddedNumber(int64(dayInt), 2)
	monthInt, _ := strconv.Atoi(date[1])
	month := PaddedNumber(int64(monthInt), 2)
	return fmt.Sprintf("%s%s%s%s%s", day, dateDelimiter, month, dateDelimiter, date[2])
}

func PaddedNumber(value int64, length int) string {
	var format string = "%0" + strconv.Itoa(length) + "d"
	return fmt.Sprintf(format, value)
}

func GetQueuePriority(companyId uint64, companyList []uint64, queuePriorityLevel uint32) uint32 {
	var queuePriority uint32

	if slices.Contains(companyList, companyId) {
		queuePriority = queuePriorityLevel
	}

	return queuePriority
}

func SetTimeLocation(timeLocation string) (*time.Location, error) {
	return time.LoadLocation(timeLocation)
}

func GetUserFromContext(ctx context.Context) (*authPB.VerifyTokenRes, *pb.UserAuthority) {
	dataClaims := ctx.Value(constant.CtxTokenKey).([]byte)
	userData := &authPB.VerifyTokenRes{}
	unmarshaller := protojson.UnmarshalOptions{
		AllowPartial: true,
	}
	unmarshaller.Unmarshal(dataClaims, userData)

	userAuthority := &pb.UserAuthority{}

	for _, authority := range userData.Authorities {
		switch authority {
		case strings.ToLower(pb.TaskStep_MAKER.String()):
			userAuthority.IsMaker = true
		case strings.ToLower(pb.TaskStep_CHECKER.String()):
			userAuthority.IsChecker = true
		case strings.ToLower(pb.TaskStep_SIGNER.String()):
			userAuthority.IsSigner = true
		case strings.ToLower(pb.TaskStep_RELEASER.String()):
			userAuthority.IsReleaser = true
		}
	}

	return userData, userAuthority
}
