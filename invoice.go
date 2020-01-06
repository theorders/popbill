package popbill


type Invoice struct {
	WriteSpecification bool   `json:"writeSpecification,omitempty" firestore:"writeSpecification"`
	ForceIssue         bool   `json:"forceIssue,omitempty" firestore:"forceIssue"`
	Memo               string `json:"memo,omitempty" firestore:"memo"`
	EmailSubject       string `json:"emailSubject,omitempty" firestore:"emailSubject"`
	DealInvoiceMgtKey  string `json:"dealInvoiceMgtKey,omitempty" firestore:"dealInvoiceMgtKey"`
	IssueType          string `json:"issueType,omitempty" firestore:"issueType"`             //!발행형태: 정발행/역발행/위수탁
	TaxType            string `json:"taxType,omitempty" firestore:"taxType"`                 //!과세형태: 과세/영세/면세
	IssueTiming        string `json:"issueTiming,omitempty" firestore:"issueTiming"`         //!발행시점: 직접발행/승인시자동발행
	ChargeDirection    string `json:"chargeDirection,omitempty" firestore:"chargeDirection"` //!과금방향: 정과금(공급자 과금)/역과금(공급받는자 과금, 역발행시에만 가능)
	SerialNum          string `json:"serialNum,omitempty" firestore:"serialNum"`
	Kwon               string `json:"kwon,omitempty" firestore:"kwon"`
	Ho                 string `json:"ho,omitempty" firestore:"ho"`
	WriteDate          string `json:"writeDate,omitempty" firestore:"writeDate"`          //!작성일자
	PurposeType        string `json:"purposeType,omitempty" firestore:"purposeType"`      //!영수/청구
	SupplyCostTotal    int64    `json:"supplyCostTotal,string" firestore:"supplyCostTotal"` //!공급가액 합계
	TaxTotal           int64    `json:"taxTotal,string" firestore:"taxTotal"`               //!세액합계
	TotalAmount        int64    `json:"totalAmount,string" firestore:"totalAmount"`         //!공급가액합계 + 세액합계
	Cash               string `json:"cash,omitempty" firestore:"cash"`
	ChkBill            string `json:"chkBill,omitempty" firestore:"chkBill"`
	Credit             string `json:"credit,omitempty" firestore:"credit"`
	Note               string `json:"note,omitempty" firestore:"note"`
	Remark1            string `json:"remark1,omitempty" firestore:"remark1"`
	Remark2            string `json:"remark2,omitempty" firestore:"remark2"`
	Remark3            string `json:"remark3,omitempty" firestore:"remark3"`
	//공급자===============================
	InvoicerMgtKey      string `json:"invoicerMgtKey,omitempty" firestore:"invoicerMgtKey"`   //*관리번호
	InvoicerCorpNum     string `json:"invoicerCorpNum,omitempty" firestore:"invoicerCorpNum"` //!사업자번호(-제외, 10자리)
	InvoicerTaxRegID    string `json:"invoicerTaxRegID,omitempty" firestore:"invoicerTaxRegID"`
	InvoicerCorpName    string `json:"invoicerCorpName,omitempty" firestore:"invoicerCorpName"` //!상호명
	InvoicerCEOName     string `json:"invoicerCEOName,omitempty" firestore:"invoicerCEOName"`   //!대표자 성명
	InvoicerAddr        string `json:"invoicerAddr,omitempty" firestore:"invoicerAddr"`
	InvoicerBizClass    string `json:"invoicerBizClass,omitempty" firestore:"invoicerBizClass"`
	InvoicerBizType     string `json:"invoicerBizType,omitempty" firestore:"invoicerBizType"`
	InvoicerContactName string `json:"invoicerContactName,omitempty" firestore:"invoicerContactName"` //!담당자명
	InvoicerDeptName    string `json:"invoicerDeptName,omitempty" firestore:"invoicerDeptName"`
	InvoicerTEL         string `json:"invoicerTEL,omitempty" firestore:"invoicerTEL"`
	InvoicerHP          string `json:"invoicerHP,omitempty" firestore:"invoicerHP"`
	InvoicerEmail       string `json:"invoicerEmail,omitempty" firestore:"invoicerEmail"`
	InvoicerSMSSendYN   bool   `json:"invoicerSMSSendYN,omitempty" firestore:"invoicerSMSSendYN"`
	//공급받는자===========================
	InvoiceeMgtKey       string `json:"invoiceeMgtKey,omitempty" firestore:"invoiceeMgtKey"`   //#문서관리번호(24자리 영문,숫자,-,_)
	InvoiceeType         string `json:"invoiceeType,omitempty" firestore:"invoiceeType"`       //!사업자/개인/외국인
	InvoiceeCorpNum      string `json:"invoiceeCorpNum,omitempty" firestore:"invoiceeCorpNum"` //!사업자 번호
	InvoiceeTaxRegID     string `json:"invoiceeTaxRegID,omitempty" firestore:"invoiceeTaxRegID"`
	InvoiceeCorpName     string `json:"invoiceeCorpName,omitempty" firestore:"invoiceeCorpName"` //!상호
	InvoiceeCEOName      string `json:"invoiceeCEOName,omitempty" firestore:"invoiceeCEOName"`   //!대표자성명
	InvoiceeAddr         string `json:"invoiceeAddr,omitempty" firestore:"invoiceeAddr"`
	InvoiceeBizClass     string `json:"invoiceeBizClass,omitempty" firestore:"invoiceeBizClass"`
	InvoiceeBizType      string `json:"invoiceeBizType,omitempty" firestore:"invoiceeBizType"`
	InvoiceeContactName1 string `json:"invoiceeContactName1,omitempty" firestore:"invoiceeContactName1"` //!담당자명
	InvoiceeDeptName1    string `json:"invoiceeDeptName1,omitempty" firestore:"invoiceeDeptName1"`
	InvoiceeTEL1         string `json:"invoiceeTEL1,omitempty" firestore:"invoiceeTEL1"`
	InvoiceeHP1          string `json:"invoiceeHP1,omitempty" firestore:"invoiceeHP1"`
	InvoiceeEmail1       string `json:"invoiceeEmail1,omitempty" firestore:"invoiceeEmail1"`
	InvoiceeContactName2 string `json:"invoiceeContactName2,omitempty" firestore:"invoiceeContactName2"`
	InvoiceeDeptName2    string `json:"invoiceeDeptName2,omitempty" firestore:"invoiceeDeptName2"`
	InvoiceeTEL2         string `json:"invoiceeTEL2,omitempty" firestore:"invoiceeTEL2"`
	InvoiceeHP2          string `json:"invoiceeHP2,omitempty" firestore:"invoiceeHP2"`
	InvoiceeEmail2       string `json:"invoiceeEmail2,omitempty" firestore:"invoiceeEmail2"`
	InvoiceeSMSSendYN    string `json:"invoiceeSMSSendYN,omitempty" firestore:"invoiceeSMSSendYN"`
	CloseDownState       bool   `json:"closeDownState,omitempty" firestore:"closeDownState"`
	CloseDownStateDate   string `json:"closeDownStateDate,omitempty" firestore:"closeDownStateDate"`
	//위수탁발행==========================
	TrusteeMgtKey      string `json:"trusteeMgtKey,omitempty" firestore:"trusteeMgtKey"`   //@문서관리번호
	TrusteeCorpNum     string `json:"trusteeCorpNum,omitempty" firestore:"trusteeCorpNum"` //@사업자번호
	TrusteeTaxRegID    string `json:"trusteeTaxRegID,omitempty" firestore:"trusteeTaxRegID"`
	TrusteeCorpName    string `json:"trusteeCorpName,omitempty" firestore:"trusteeCorpName"` //@상호
	TrusteeCEOName     string `json:"trusteeCEOName,omitempty" firestore:"trusteeCEOName"`   //@대표자성명
	TrusteeAddr        string `json:"trusteeAddr,omitempty" firestore:"trusteeAddr"`
	TrusteeBizClass    string `json:"trusteeBizClass,omitempty" firestore:"trusteeBizClass"`
	TrusteeBizType     string `json:"trusteeBizType,omitempty" firestore:"trusteeBizType"`
	TrusteeContactName string `json:"trusteeContactName,omitempty" firestore:"trusteeContactName"` //@담당자명
	TrusteeDeptName    string `json:"trusteeDeptName,omitempty" firestore:"trusteeDeptName"`
	TrusteeTEL         string `json:"trusteeTEL,omitempty" firestore:"trusteeTEL"`
	TrusteeHP          string `json:"trusteeHP,omitempty" firestore:"trusteeHP"`
	TrusteeEmail       string `json:"trusteeEmail,omitempty" firestore:"trusteeEmail"`
	TrusteeSMSSendYN   bool   `json:"trusteeSMSSendYN,omitempty" firestore:"trusteeSMSSendYN"`
	//수정세금계산서=======================
	//1 : 기재사항 착오정정
	//2 : 공급가액 변동
	//3 : 환입
	//4 : 계약의 해지
	//5 : 내국신용장 사후개설
	//6 : 착오에 의한 이중발행
	ModifyCode            int64  `json:"modifyCode,omitempty" firestore:"modifyCode"`                       //수정사유코드
	OriginalTaxinvoiceKey string `json:"originalTaxinvoiceKey,omitempty" firestore:"originalTaxinvoiceKey"` //원본팝빌관리번호

	DetailList     []Detail  `json:"detailList,omitempty" firestore:"detailList"`
	AddContactList []Contact `json:"addContactList,omitempty" firestore:"addContactList"`
}

type Detail struct {
	SerialNum  int    `json:"serialNum,omitempty" firestore:"serialNum"`
	PurchaseDT string `json:"purchaseDT,omitempty" firestore:"purchaseDT"`
	ItemName   string `json:"itemName,omitempty" firestore:"itemName"`
	Spec       string `json:"spec,omitempty" firestore:"spec"`
	Qty        string `json:"qty,omitempty" firestore:"qty"`
	UnitCost   int64    `json:"unitCost,string" firestore:"unitCost"`
	SupplyCost int64    `json:"supplyCost,string" firestore:"supplyCost"`
	Tax        int64    `json:"tax,string" firestore:"tax"`
	Remark     string `json:"remark,omitempty" firestore:"remark"`
}

type Contact struct {
	SerialNum   int    `json:"serialNum,string" firestore:"serialNum"`        //!일련번호
	ContactName string `json:"contactName,omitempty" firestore:"contactName"` //!담당자명
	Email       string `json:"email,omitempty" firestore:"email"`             //!이메일주소
}


