package entity

import (
	"testing"
	"time"

	"github.com/asaskevich/govalidator"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

type BorrowBook struct {
	gorm.Model
	Borb_Day       time.Time `valid:"required~วันยืมหนังสือต้องเป็นปัจจุบัน, Present~วันยืมหนังสือต้องเป็นปัจจุบัน"`
	Return_Day     time.Time `valid:"required~วันคืนหนังสือต้องไม่เป็นอดีต, Past~วันคืนหนังสือต้องไม่เป็นอดีต"`
	Color_Bar      string
	Book_Frequency int `valid:"required~ต้องเป็นตัวเลข 1-1000, range(1|1000)~ต้องเป็นตัวเลข 1-1000"`
}

func TestAllBorrowBookCorrect(t *testing.T) {
	g := NewGomegaWithT(t)

	borrowbook := BorrowBook{
		Borb_Day:       time.Now(),
		Return_Day:     time.Now(),
		Color_Bar:      "สีเเดง",
		Book_Frequency: 1,
	}
	ok, err := govalidator.ValidateStruct(borrowbook)
	g.Expect(ok).To(BeTrue())
	g.Expect(err).To(BeNil())
}

func TestBorb_DayMustBePresent(t *testing.T) {
	g := NewGomegaWithT(t)

	fixture := []time.Time{
		time.Now().Add(+24 * time.Hour),
		time.Now().Add(-24 * time.Hour),
	}

	for _, borbday := range fixture {
		borrowbook := BorrowBook{
			Borb_Day:       borbday, //ผิด
			Return_Day:     time.Now(),
			Color_Bar:      "สีเเดง",
			Book_Frequency: 1,
		}
		ok, err := govalidator.ValidateStruct(borrowbook)
		g.Expect(ok).ToNot(BeTrue())
		g.Expect(err).ToNot(BeNil())
		g.Expect(err.Error()).To(Equal("วันยืมหนังสือต้องเป็นปัจจุบัน"))
	}
}

func TestReturn_DayMustNotBePast(t *testing.T) {
	g := NewGomegaWithT(t)

	borrowbook := BorrowBook{
		Borb_Day:       time.Now(),
		Return_Day:     time.Now().Add(-24 * time.Hour), //ผิด
		Color_Bar:      "สีเเดง",
		Book_Frequency: 1,
	}
	ok, err := govalidator.ValidateStruct(borrowbook)
	g.Expect(ok).ToNot(BeTrue())
	g.Expect(err).ToNot(BeNil())
	g.Expect(err.Error()).To(Equal("วันคืนหนังสือต้องไม่เป็นอดีต"))
}

// func TestColor_Bar(t *testing.T) {
// 	g := NewGomegaWithT(t)

// 	borrowbook := BorrowBook{
// 		Borb_Day:       time.Now(),
// 		Return_Day:     time.Now(),
// 		Color_Bar:      "*฿", //ผิด
// 		Book_Frequency: 1,
// 	}
// 	ok, err := govalidator.ValidateStruct(borrowbook)
// 	g.Expect(ok).ToNot(BeTrue())
// 	g.Expect(err).ToNot(BeNil())
// 	g.Expect(err.Error()).To(Equal("เเถบสีหนังสือต้องไม่เป็นอักขระพิเศษ"))
// }

func TestBook_Frequency(t *testing.T) {
	g := NewGomegaWithT(t)

	fixture := []int{
		-1, 0, 1001,
	}

	for _, book := range fixture {
		borrowbook := BorrowBook{
			Borb_Day:       time.Now(), //ผิด
			Return_Day:     time.Now(),
			Color_Bar:      "สีเเดง",
			Book_Frequency: book,
		}
		ok, err := govalidator.ValidateStruct(borrowbook)
		g.Expect(ok).ToNot(BeTrue())
		g.Expect(err).ToNot(BeNil())
		g.Expect(err.Error()).To(Equal("ต้องเป็นตัวเลข 1-1000"))
	}
}

func init() {
	govalidator.CustomTypeTagMap.Set("Past", func(i interface{}, context interface{}) bool {
		t := i.(time.Time)
		return t.After(time.Now().Add(time.Minute*-2)) || t.Equal(time.Now())
	})

	govalidator.CustomTypeTagMap.Set("Present", func(i interface{}, context interface{}) bool {
		t := i.(time.Time)
		return t.After(time.Now().Add(-2*time.Minute)) && t.Before(time.Now().Add(+2*time.Minute))
	})

	govalidator.CustomTypeTagMap.Set("Future", func(i interface{}, context interface{}) bool {
		t := i.(time.Time)
		return t.Before(time.Now().Add(+2*time.Minute)) || t.Equal(time.Now())
	})
}

// func init() {
// 	govalidator.CustomTypeTagMap.Set("cha_valid", govalidator.CustomTypeValidator(func(i interface{}, context interface{}) bool {
// 		s, ok := i.(string)
// 		if !ok {
// 			return false
// 		}
// 		match, _ := regexp.MatchString("^[ก-ฮa-zA-Z\\s]+$", s)
// 		return match
// 	}))
// }
