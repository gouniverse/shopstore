package shopstore

import (
	"database/sql"
	"os"
	"strings"
	"testing"

	"github.com/gouniverse/sb"
	_ "modernc.org/sqlite"
)

func initDB(filepath string) *sql.DB {
	os.Remove(filepath) // remove database
	dsn := filepath + "?parseTime=true"
	db, err := sql.Open("sqlite", dsn)

	if err != nil {
		panic(err)
	}

	return db
}

func TestStoreDiscountCreate(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		DiscountTableName:  "shop_discount_create",
		OrderTableName:     "shop_order_create",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	discount := NewDiscount().
		SetStatus(DISCOUNT_STATUS_DRAFT).
		SetTitle("DISCOUNT_TITLE")

	err = store.DiscountCreate(discount)

	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreDiscountFindByID(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		DiscountTableName:  "shop_discount_find_by_id",
		OrderTableName:     "shop_order_find_by_id",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	discount := NewDiscount().
		SetStatus(DISCOUNT_STATUS_DRAFT).
		SetTitle("DISCOUNT_TITLE").
		SetDescription("DISCOUNT_DESCRIPTION").
		SetType(DISCOUNT_TYPE_AMOUNT).
		SetAmount(19.99).
		SetStartsAt(`2022-01-01 00:00:00`).
		SetEndsAt(`2022-01-01 23:59:59`)

	err = store.DiscountCreate(discount)

	if err != nil {
		t.Fatal("unexpected error:", err)
		return
	}

	discountFound, errFind := store.DiscountFindByID(discount.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
		return
	}

	if discountFound == nil {
		t.Fatal("Exam MUST NOT be nil")
		return
	}

	if discountFound.Title() != "DISCOUNT_TITLE" {
		t.Fatal("Exam title MUST BE 'DISCOUNT_TITLE', found: ", discountFound.Title())
		return
	}

	if discountFound.Description() != "DISCOUNT_DESCRIPTION" {
		t.Fatal("Exam description MUST BE 'DISCOUNT_DESCRIPTION', found: ", discountFound.Description())
	}

	if discountFound.Status() != DISCOUNT_STATUS_DRAFT {
		t.Fatal("Exam status MUST BE 'draft', found: ", discountFound.Status())
		return
	}

	if discountFound.Type() != DISCOUNT_TYPE_AMOUNT {
		t.Fatal("Exam type MUST BE 'amount', found: ", discountFound.Type())
	}

	if discountFound.Type() != DISCOUNT_TYPE_AMOUNT {
		t.Fatal("Exam type MUST BE 'amount', found: ", discountFound.Type())
	}

	if discountFound.Amount() != 19.9900 {
		t.Fatal("Exam price MUST BE '19.9900', found: ", discountFound.Amount())
		return
	}

	if discountFound.StartsAt() != "2022-01-01 00:00:00 +0000 UTC" {
		t.Fatal("Exam start date MUST BE '2022-01-01 00:00:00', found: ", discountFound.StartsAt())
	}

	if discountFound.EndsAt() != "2022-01-01 23:59:59 +0000 UTC" {
		t.Fatal("Exam end date MUST BE '2022-01-01 23:59:59', found: ", discountFound.EndsAt())
	}

	// if examFound.Memo() != "test memo" {
	// 	t.Fatal("Exam memo MUST BE 'test memo', found: ", examFound.Memo())
	// }

	if !strings.Contains(discountFound.DeletedAt(), sb.NULL_DATETIME) {
		t.Fatal("Exam MUST NOT be soft deleted", discountFound.DeletedAt())
		return
	}
}

// func TestExamServiceExamSoftDelete(t *testing.T) {
// 	config.TestsConfigureAndInitialize()
// 	Initialize()

// 	exam := NewExam().
// 		SetStatus(EXAM_STATUS_DRAFT).
// 		SetTitle("EXAM01_ID").
// 		SetPrice("19.99").
// 		SetQuestionsNumber("30").
// 		SetMinutesToComplete("60")

// 	err := NewExamService().ExamCreate(exam)
// 	if err != nil {
// 		t.Fatal("unexpected error:", err)
// 	}

// 	err = NewExamService().ExamSoftDeleteByID(exam.ID())

// 	if err != nil {
// 		t.Fatal("unexpected error:", err)
// 	}

// 	if exam.DeletedAt() != sb.NULL_DATETIME {
// 		t.Fatal("Exam MUST NOT be soft deleted")
// 	}

// 	examFound, errFind := NewExamService().ExamFindByID(exam.ID())

// 	if errFind != nil {
// 		t.Fatal("unexpected error:", errFind)
// 	}

// 	if examFound != nil {
// 		t.Fatal("Exam MUST be nil")
// 	}

// 	examFindWithDeleted, err := NewExamService().ExamList(ExamQueryOptions{
// 		ID:          exam.ID(),
// 		Limit:       1,
// 		WithDeleted: true,
// 	})

// 	if err != nil {
// 		t.Fatal("unexpected error:", err)
// 	}

// 	if len(examFindWithDeleted) == 0 {
// 		t.Fatal("Exam MUST be soft deleted")
// 	}

// 	if strings.Contains(examFindWithDeleted[0].DeletedAt(), sb.NULL_DATETIME) {
// 		t.Fatal("Exam MUST be soft deleted", examFound.DeletedAt())
// 	}

// }

// func TestExamServiceQuestionCreate(t *testing.T) {
// 	config.TestsConfigureAndInitialize()
// 	Initialize()
// 	question := NewExamQuestion().
// 		SetStatus(EXAM_QUESTION_STATUS_DRAFT).
// 		SetType(EXAM_QUESTION_TYPE_CHOICE_SINGLE).
// 		SetTitle("QUESTION_TITLE")

// 	err := NewExamService().QuestionCreate(question)
// 	if err != nil {
// 		t.Fatal("unexpected error:", err)
// 	}
// }

// func TestExamServiceQuestionFindByID(t *testing.T) {
// 	config.TestsConfigureAndInitialize()
// 	Initialize()

// 	question := NewExamQuestion().
// 		SetExamID("EXAM01").
// 		SetStatus(EXAM_QUESTION_STATUS_DRAFT).
// 		SetType(EXAM_QUESTION_TYPE_CHOICE_SINGLE).
// 		SetTitle("QUESTION_TITLE").
// 		SetDetails("QUESTION_DETAILS").
// 		SetMemo("test memo")

// 	err := NewExamService().QuestionCreate(question)
// 	if err != nil {
// 		t.Fatal("unexpected error:", err)
// 	}

// 	questionFound, errFind := NewExamService().QuestionFindByID(question.ID())

// 	if errFind != nil {
// 		t.Fatal("unexpected error:", errFind)
// 	}

// 	if questionFound == nil {
// 		t.Fatal("Exam MUST NOT be nil")
// 	}

// 	if questionFound.ExamID() != "EXAM01" {
// 		t.Fatal("Question exam_id MUST BE 'EXAM01', found: ", questionFound.ExamID())
// 	}

// 	if questionFound.Title() != "QUESTION_TITLE" {
// 		t.Fatal("Question title MUST BE 'QUESTION_TITLE', found: ", questionFound.Title())
// 	}

// 	if questionFound.Status() != EXAM_QUESTION_STATUS_DRAFT {
// 		t.Fatal("Question status MUST BE 'draft', found: ", questionFound.Status())
// 	}

// 	if questionFound.Details() != "QUESTION_DETAILS" {
// 		t.Fatal("Question details MUST BE 'QUESTION_DETAILS', found: ", questionFound.Details())
// 	}

// 	if questionFound.Type() != EXAM_QUESTION_TYPE_CHOICE_SINGLE {
// 		t.Fatal("Question type MUST BE 'choice_single', found: ", questionFound.Type())
// 	}

// 	if questionFound.Memo() != "test memo" {
// 		t.Fatal("Exam memo MUST BE 'test memo', found: ", questionFound.Memo())
// 	}

// 	if !strings.Contains(questionFound.DeletedAt(), sb.NULL_DATETIME) {
// 		t.Fatal("Question MUST NOT be soft deleted", questionFound.DeletedAt())
// 	}
// }

// func TestExamServiceQuestionSoftDelete(t *testing.T) {
// 	config.TestsConfigureAndInitialize()
// 	Initialize()

// 	question := NewExamQuestion().
// 		SetStatus(EXAM_QUESTION_STATUS_DRAFT).
// 		SetType(EXAM_QUESTION_TYPE_CHOICE_SINGLE).
// 		SetTitle("QUESTION_TITLE")

// 	err := NewExamService().QuestionCreate(question)
// 	if err != nil {
// 		t.Fatal("unexpected error:", err)
// 	}

// 	err = NewExamService().QuestionSoftDeleteByID(question.ID())

// 	if err != nil {
// 		t.Fatal("unexpected error:", err)
// 	}

// 	if question.DeletedAt() != sb.NULL_DATETIME {
// 		t.Fatal("Question MUST NOT be soft deleted")
// 	}

// 	questionFound, errFind := NewExamService().QuestionFindByID(question.ID())

// 	if errFind != nil {
// 		t.Fatal("unexpected error:", errFind)
// 	}

// 	if questionFound != nil {
// 		t.Fatal("Question MUST be nil")
// 	}

// 	questionFindWithDeleted, err := NewExamService().QuestionList(ExamQuestionQueryOptions{
// 		ID:          question.ID(),
// 		Limit:       1,
// 		WithDeleted: true,
// 	})

// 	if err != nil {
// 		t.Fatal("unexpected error:", err)
// 	}

// 	if len(questionFindWithDeleted) == 0 {
// 		t.Fatal("Exam MUST be soft deleted")
// 	}

// 	if strings.Contains(questionFindWithDeleted[0].DeletedAt(), sb.NULL_DATETIME) {
// 		t.Fatal("Question MUST be soft deleted", questionFound.DeletedAt())
// 	}

// }

// func TestExamServiceOptionCreate(t *testing.T) {
// 	config.TestsConfigureAndInitialize()
// 	Initialize()
// 	option := NewExamQuestionOption().
// 		SetStatus(EXAM_QUESTION_OPTION_STATUS_DRAFT).
// 		SetTitle("OPTION_TITLE")

// 	err := NewExamService().OptionCreate(option)
// 	if err != nil {
// 		t.Fatal("unexpected error:", err)
// 	}
// }

// func TestExamServiceOptionFindByID(t *testing.T) {
// 	config.TestsConfigureAndInitialize()
// 	Initialize()

// 	option := NewExamQuestionOption().
// 		SetStatus(EXAM_QUESTION_OPTION_STATUS_DRAFT).
// 		SetQuestionID("QUESTION01").
// 		SetTitle("OPTION_TITLE").
// 		SetDetails("OPTION_DETAILS").
// 		SetMemo("test memo")

// 	err := NewExamService().OptionCreate(option)
// 	if err != nil {
// 		t.Fatal("unexpected error:", err)
// 	}

// 	optionFound, errFind := NewExamService().OptionFindByID(option.ID())

// 	if errFind != nil {
// 		t.Fatal("unexpected error:", errFind)
// 	}

// 	if optionFound == nil {
// 		t.Fatal("Option MUST NOT be nil")
// 	}

// 	if optionFound.QuestionID() != "QUESTION01" {
// 		t.Fatal("Option question_id MUST BE 'QUESTION01', found: ", optionFound.QuestionID())
// 	}

// 	if optionFound.Title() != "OPTION_TITLE" {
// 		t.Fatal("Option title MUST BE 'OPTION_TITLE', found: ", optionFound.Title())
// 	}

// 	if optionFound.Status() != EXAM_STATUS_DRAFT {
// 		t.Fatal("Option status MUST BE 'draft', found: ", optionFound.Status())
// 	}

// 	if optionFound.Details() != "OPTION_DETAILS" {
// 		t.Fatal("Option details MUST BE 'OPTION_DETAILS', found: ", optionFound.Details())
// 	}

// 	if optionFound.Memo() != "test memo" {
// 		t.Fatal("Exam memo MUST BE 'test memo', found: ", optionFound.Memo())
// 	}

// 	if !strings.Contains(optionFound.DeletedAt(), sb.NULL_DATETIME) {
// 		t.Fatal("Option MUST NOT be soft deleted", optionFound.DeletedAt())
// 	}
// }

// func TestExamServiceOptionSoftDelete(t *testing.T) {
// 	config.TestsConfigureAndInitialize()
// 	Initialize()

// 	option := NewExamQuestionOption().
// 		SetStatus(EXAM_QUESTION_OPTION_STATUS_DRAFT).
// 		SetTitle("OPTION_TITLE").
// 		SetDetails("OPTION_DETAILS").
// 		SetMemo("test memo")

// 	err := NewExamService().OptionCreate(option)
// 	if err != nil {
// 		t.Fatal("unexpected error:", err)
// 	}

// 	err = NewExamService().OptionSoftDeleteByID(option.ID())

// 	if err != nil {
// 		t.Fatal("unexpected error:", err)
// 	}

// 	if option.DeletedAt() != sb.NULL_DATETIME {
// 		t.Fatal("Exam MUST NOT be soft deleted")
// 	}

// 	optionFound, errFind := NewExamService().OptionFindByID(option.ID())

// 	if errFind != nil {
// 		t.Fatal("unexpected error:", errFind)
// 	}

// 	if optionFound != nil {
// 		t.Fatal("Exam MUST be nil")
// 	}

// 	optionFindWithDeleted, err := NewExamService().OptionList(ExamQuestionOptionQueryOptions{
// 		ID:          option.ID(),
// 		Limit:       1,
// 		WithDeleted: true,
// 	})

// 	if err != nil {
// 		t.Fatal("unexpected error:", err)
// 	}

// 	if len(optionFindWithDeleted) == 0 {
// 		t.Fatal("Exam MUST be soft deleted")
// 	}

// 	if strings.Contains(optionFindWithDeleted[0].DeletedAt(), sb.NULL_DATETIME) {
// 		t.Fatal("Exam MUST be soft deleted", optionFound.DeletedAt())
// 	}

// }

func TestStoreOderCreate(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		DiscountTableName:  "shop_discount_create",
		OrderTableName:     "shop_order_create",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	order := NewOrder().
		SetStatus(ORDER_STATUS_PENDING).
		SetUserID("USER)1_ID").
		SetExamID("EXAM01_ID").
		SetQuantity(1).
		SetPrice(19.99)

	err = store.OrderCreate(order)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
}

func TestStoreOrderFindByID(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		DiscountTableName:  "shop_discount_create",
		OrderTableName:     "shop_order_create",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	order := NewOrder().
		SetStatus(ORDER_STATUS_PENDING).
		SetUserID("USER01_ID").
		SetExamID("EXAM01_ID").
		SetQuantity(1).
		SetPrice(19.99).
		SetMemo("test memo")

	err = store.OrderCreate(order)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	orderFound, errFind := store.OrderFindByID(order.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if orderFound == nil {
		t.Fatal("Order MUST NOT be nil")
	}

	if orderFound.UserID() != "USER01_ID" {
		t.Fatal("Order user id MUST BE 'USER01_ID', found: ", orderFound.UserID())
	}

	if orderFound.ExamID() != "EXAM01_ID" {
		t.Fatal("Order exam id MUST BE 'EXAM01_ID', found: ", orderFound.ExamID())
	}

	if orderFound.Status() != ORDER_STATUS_PENDING {
		t.Fatal("Order status MUST BE 'pending', found: ", orderFound.Status())
	}

	if orderFound.Quantity() != "1" {
		t.Fatal("Order quantity MUST BE '1', found: ", orderFound.Quantity())
	}

	if orderFound.Price() != "19.9900" {
		t.Fatal("Order price MUST BE '19.9900', found: ", orderFound.Price())
	}

	if orderFound.Memo() != "test memo" {
		t.Fatal("Order memo MUST BE 'test memo', found: ", orderFound.Memo())
	}

	if !strings.Contains(orderFound.DeletedAt(), sb.NULL_DATETIME) {
		t.Fatal("Order MUST NOT be soft deleted", orderFound.DeletedAt())
	}
}

func TestStoreOrderSoftDelete(t *testing.T) {
	db := initDB(":memory:")

	store, err := NewStore(NewStoreOptions{
		DB:                 db,
		DiscountTableName:  "shop_discount_create",
		OrderTableName:     "shop_order_create",
		AutomigrateEnabled: true,
	})

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if store == nil {
		t.Fatal("unexpected nil store")
	}

	order := NewOrder().
		SetStatus(ORDER_STATUS_PENDING).
		SetUserID("USER01_ID").
		SetExamID("EXAM01_ID").
		SetQuantity(1).
		SetPrice(19.99)

	err = store.OrderCreate(order)
	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	err = store.OrderSoftDeleteByID(order.ID())

	if err != nil {
		t.Fatal("unexpected error:", err)
	}

	if order.DeletedAt() != sb.NULL_DATETIME {
		t.Fatal("Order MUST NOT be soft deleted")
	}

	orderFound, errFind := store.OrderFindByID(order.ID())

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if orderFound != nil {
		t.Fatal("Order MUST be nil")
	}

	orderFindWithDeleted, errFind := store.OrderList(OrderQueryOptions{
		ID:          order.ID(),
		Limit:       1,
		WithDeleted: true,
	})

	if errFind != nil {
		t.Fatal("unexpected error:", errFind)
	}

	if len(orderFindWithDeleted) < 1 {
		t.Fatal("Order list MUST NOT be empty")
		return
	}

	if strings.Contains(orderFindWithDeleted[0].DeletedAt(), sb.NULL_DATETIME) {
		t.Fatal("Order MUST be soft deleted", orderFound.DeletedAt())
	}

}
