package frame

type Type byte

const (
	TypeAuthRequest          = 'R'
	TypeParameterStatus      = 'S'
	TypeBackendKeyData       = 'K'
	TypeReadyForQuery        = 'Z'
	TypeParseCompletion      = '1'
	TypeParameterDescription = 't'
	TypeRowDescription       = 'T'
	TypeDataRow              = 'D'
	TypeCommandCompletion    = 'C'
	TypeError                = 'E'
	TypeNotification         = 'A'
	TypeCloseComplete        = '3'
	TypeEmptyQueryResponse   = 'I'
	TypeNoData               = 'n'
	TypeNoticeResponse       = 'N'
)

type TransactionStatus byte

const (
	TransactionStatusNoReady             = 'N'
	TransactionStatusIdle                = 'I'
	TransactionStatusIdleInTransaction   = 'T'
	TransactionStatusInFailedTransaction = 'E'
)

type PgType uint32

const (
	PgTypeBool             = 16
	PgTypeBytea            = 17
	PgTypeChar             = 18
	PgTypeName             = 19
	PgTypeInt8             = 20
	PgTypeInt2             = 21
	PgTypeInt2vector       = 22
	PgTypeInt4             = 23
	PgTypeRegproc          = 24
	PgTypeText             = 25
	PgTypeOid              = 26
	PgTypeTid              = 27
	PgTypeXid              = 28
	PgTypeCid              = 29
	PgTypeOidvector        = 30
	PgTypePgDdlCommand     = 32
	PgTypePgType           = 71
	PgTypePgAttribute      = 75
	PgTypePgProc           = 81
	PgTypePgClass          = 83
	PgTypeJson             = 114
	PgTypeXml              = 142
	PgTypeArrXml           = 143
	PgTypePgNodeTree       = 194
	PgTypeArrJson          = 199
	PgTypeSmgr             = 210
	PgTypeIndexAmHandler   = 325
	PgTypePoint            = 600
	PgTypeLseg             = 601
	PgTypePath             = 602
	PgTypeBox              = 603
	PgTypePolygon          = 604
	PgTypeLine             = 628
	PgTypeArrLine          = 629
	PgTypeCidr             = 650
	PgTypeArrCidr          = 651
	PgTypeFloat4           = 700
	PgTypeFloat8           = 701
	PgTypeAbstime          = 702
	PgTypeReltime          = 703
	PgTypeTinterval        = 704
	PgTypeUnknown          = 705
	PgTypeCircle           = 718
	PgTypeArrCircle        = 719
	PgTypeMacaddr8         = 774
	PgTypeArrMacaddr8      = 775
	PgTypeMoney            = 790
	PgTypeArrMoney         = 791
	PgTypeMacaddr          = 829
	PgTypeInet             = 869
	PgTypeArrBool          = 1000
	PgTypeArrBytea         = 1001
	PgTypeArrChar          = 1002
	PgTypeArrName          = 1003
	PgTypeArrInt2          = 1005
	PgTypeArrInt2vector    = 1006
	PgTypeArrInt4          = 1007
	PgTypeArrRegproc       = 1008
	PgTypeArrText          = 1009
	PgTypeArrTid           = 1010
	PgTypeArrXid           = 1011
	PgTypeArrCid           = 1012
	PgTypeArrOidvector     = 1013
	PgTypeArrBpchar        = 1014
	PgTypeArrVarchar       = 1015
	PgTypeArrInt8          = 1016
	PgTypeArrPoint         = 1017
	PgTypeArrLseg          = 1018
	PgTypeArrPath          = 1019
	PgTypeArrBox           = 1020
	PgTypeArrFloat4        = 1021
	PgTypeArrFloat8        = 1022
	PgTypeArrAbstime       = 1023
	PgTypeArrReltime       = 1024
	PgTypeArrTinterval     = 1025
	PgTypeArrPolygon       = 1027
	PgTypeArrOid           = 1028
	PgTypeAclitem          = 1033
	PgTypeArrAclitem       = 1034
	PgTypeArrMacaddr       = 1040
	PgTypeArrInet          = 1041
	PgTypeBpchar           = 1042
	PgTypeVarchar          = 1043
	PgTypeDate             = 1082
	PgTypeTime             = 1083
	PgTypeTimestamp        = 1114
	PgTypeArrTimestamp     = 1115
	PgTypeArrDate          = 1182
	PgTypeArrTime          = 1183
	PgTypeTimestamptz      = 1184
	PgTypeArrTimestamptz   = 1185
	PgTypeInterval         = 1186
	PgTypeArrInterval      = 1187
	PgTypeArrNumeric       = 1231
	PgTypePgDatabase       = 1248
	PgTypeArrCstring       = 1263
	PgTypeTimetz           = 1266
	PgTypeArrTimetz        = 1270
	PgTypeBit              = 1560
	PgTypeArrBit           = 1561
	PgTypeVarbit           = 1562
	PgTypeArrVarbit        = 1563
	PgTypeNumeric          = 1700
	PgTypeRefcursor        = 1790
	PgTypeArrRefcursor     = 2201
	PgTypeRegprocedure     = 2202
	PgTypeRegoper          = 2203
	PgTypeRegoperator      = 2204
	PgTypeRegclass         = 2205
	PgTypeRegtype          = 2206
	PgTypeArrRegprocedure  = 2207
	PgTypeArrRegoper       = 2208
	PgTypeArrRegoperator   = 2209
	PgTypeArrRegclass      = 2210
	PgTypeArrRegtype       = 2211
	PgTypeRecord           = 2249
	PgTypeCstring          = 2275
	PgTypeAny              = 2276
	PgTypeAnyarray         = 2277
	PgTypeVoid             = 2278
	PgTypeTrigger          = 2279
	PgTypeLanguageHandler  = 2280
	PgTypeInternal         = 2281
	PgTypeOpaque           = 2282
	PgTypeAnyelement       = 2283
	PgTypeArrRecord        = 2287
	PgTypeAnynonarray      = 2776
	PgTypePgAuthid         = 2842
	PgTypePgAuthMembers    = 2843
	PgTypeArrTxidSnapshot  = 2949
	PgTypeUuid             = 2950
	PgTypeArrUuid          = 2951
	PgTypeTxidSnapshot     = 2970
	PgTypeFdwHandler       = 3115
	PgTypePgLsn            = 3220
	PgTypeArrPgLsn         = 3221
	PgTypeTsmHandler       = 3310
	PgTypePgNdistinct      = 3361
	PgTypePgDependencies   = 3402
	PgTypeAnyenum          = 3500
	PgTypeTsvector         = 3614
	PgTypeTsquery          = 3615
	PgTypeGtsvector        = 3642
	PgTypeArrTsvector      = 3643
	PgTypeArrGtsvector     = 3644
	PgTypeArrTsquery       = 3645
	PgTypeRegconfig        = 3734
	PgTypeArrRegconfig     = 3735
	PgTypeRegdictionary    = 3769
	PgTypeArrRegdictionary = 3770
	PgTypeJsonb            = 3802
	PgTypeArrJsonb         = 3807
	PgTypeAnyrange         = 3831
	PgTypeEventTrigger     = 3838
	PgTypeInt4range        = 3904
	PgTypeArrInt4range     = 3905
	PgTypeNumrange         = 3906
	PgTypeArrNumrange      = 3907
	PgTypeTsrange          = 3908
	PgTypeArrTsrange       = 3909
	PgTypeTstzrange        = 3910
	PgTypeArrTstzrange     = 3911
	PgTypeDaterange        = 3912
	PgTypeArrDaterange     = 3913
	PgTypeInt8range        = 3926
	PgTypeArrInt8range     = 3927
	PgTypePgShseclabel     = 4066
	PgTypeRegnamespace     = 4089
	PgTypeArrRegnamespace  = 4090
	PgTypeRegrole          = 4096
	PgTypeArrRegrole       = 4097
	PgTypePgSubscription   = 6101
)

var PgTypeMap = map[PgType]string{
	PgTypeBool:             "PgTypeBool",
	PgTypeBytea:            "PgTypeBytea",
	PgTypeChar:             "PgTypeChar",
	PgTypeName:             "PgTypeName",
	PgTypeInt8:             "PgTypeInt8",
	PgTypeInt2:             "PgTypeInt2",
	PgTypeInt2vector:       "PgTypeInt2vector",
	PgTypeInt4:             "PgTypeInt4",
	PgTypeRegproc:          "PgTypeRegproc",
	PgTypeText:             "PgTypeText",
	PgTypeOid:              "PgTypeOid",
	PgTypeTid:              "PgTypeTid",
	PgTypeXid:              "PgTypeXid",
	PgTypeCid:              "PgTypeCid",
	PgTypeOidvector:        "PgTypeOidvector",
	PgTypePgDdlCommand:     "PgTypePgDdlCommand",
	PgTypePgType:           "PgTypePgType",
	PgTypePgAttribute:      "PgTypePgAttribute",
	PgTypePgProc:           "PgTypePgProc",
	PgTypePgClass:          "PgTypePgClass",
	PgTypeJson:             "PgTypeJson",
	PgTypeXml:              "PgTypeXml",
	PgTypeArrXml:           "PgTypeArrXml",
	PgTypePgNodeTree:       "PgTypePgNodeTree",
	PgTypeArrJson:          "PgTypeArrJson",
	PgTypeSmgr:             "PgTypeSmgr",
	PgTypeIndexAmHandler:   "PgTypeIndexAmHandler",
	PgTypePoint:            "PgTypePoint",
	PgTypeLseg:             "PgTypeLseg",
	PgTypePath:             "PgTypePath",
	PgTypeBox:              "PgTypeBox",
	PgTypePolygon:          "PgTypePolygon",
	PgTypeLine:             "PgTypeLine",
	PgTypeArrLine:          "PgTypeArrLine",
	PgTypeCidr:             "PgTypeCidr",
	PgTypeArrCidr:          "PgTypeArrCidr",
	PgTypeFloat4:           "PgTypeFloat4",
	PgTypeFloat8:           "PgTypeFloat8",
	PgTypeAbstime:          "PgTypeAbstime",
	PgTypeReltime:          "PgTypeReltime",
	PgTypeTinterval:        "PgTypeTinterval",
	PgTypeUnknown:          "PgTypeUnknown",
	PgTypeCircle:           "PgTypeCircle",
	PgTypeArrCircle:        "PgTypeArrCircle",
	PgTypeMacaddr8:         "PgTypeMacaddr8",
	PgTypeArrMacaddr8:      "PgTypeArrMacaddr8",
	PgTypeMoney:            "PgTypeMoney",
	PgTypeArrMoney:         "PgTypeArrMoney",
	PgTypeMacaddr:          "PgTypeMacaddr",
	PgTypeInet:             "PgTypeInet",
	PgTypeArrBool:          "PgTypeArrBool",
	PgTypeArrBytea:         "PgTypeArrBytea",
	PgTypeArrChar:          "PgTypeArrChar",
	PgTypeArrName:          "PgTypeArrName",
	PgTypeArrInt2:          "PgTypeArrInt2",
	PgTypeArrInt2vector:    "PgTypeArrInt2vector",
	PgTypeArrInt4:          "PgTypeArrInt4",
	PgTypeArrRegproc:       "PgTypeArrRegproc",
	PgTypeArrText:          "PgTypeArrText",
	PgTypeArrTid:           "PgTypeArrTid",
	PgTypeArrXid:           "PgTypeArrXid",
	PgTypeArrCid:           "PgTypeArrCid",
	PgTypeArrOidvector:     "PgTypeArrOidvector",
	PgTypeArrBpchar:        "PgTypeArrBpchar",
	PgTypeArrVarchar:       "PgTypeArrVarchar",
	PgTypeArrInt8:          "PgTypeArrInt8",
	PgTypeArrPoint:         "PgTypeArrPoint",
	PgTypeArrLseg:          "PgTypeArrLseg",
	PgTypeArrPath:          "PgTypeArrPath",
	PgTypeArrBox:           "PgTypeArrBox",
	PgTypeArrFloat4:        "PgTypeArrFloat4",
	PgTypeArrFloat8:        "PgTypeArrFloat8",
	PgTypeArrAbstime:       "PgTypeArrAbstime",
	PgTypeArrReltime:       "PgTypeArrReltime",
	PgTypeArrTinterval:     "PgTypeArrTinterval",
	PgTypeArrPolygon:       "PgTypeArrPolygon",
	PgTypeArrOid:           "PgTypeArrOid",
	PgTypeAclitem:          "PgTypeAclitem",
	PgTypeArrAclitem:       "PgTypeArrAclitem",
	PgTypeArrMacaddr:       "PgTypeArrMacaddr",
	PgTypeArrInet:          "PgTypeArrInet",
	PgTypeBpchar:           "PgTypeBpchar",
	PgTypeVarchar:          "PgTypeVarchar",
	PgTypeDate:             "PgTypeDate",
	PgTypeTime:             "PgTypeTime",
	PgTypeTimestamp:        "PgTypeTimestamp",
	PgTypeArrTimestamp:     "PgTypeArrTimestamp",
	PgTypeArrDate:          "PgTypeArrDate",
	PgTypeArrTime:          "PgTypeArrTime",
	PgTypeTimestamptz:      "PgTypeTimestamptz",
	PgTypeArrTimestamptz:   "PgTypeArrTimestamptz",
	PgTypeInterval:         "PgTypeInterval",
	PgTypeArrInterval:      "PgTypeArrInterval",
	PgTypeArrNumeric:       "PgTypeArrNumeric",
	PgTypePgDatabase:       "PgTypePgDatabase",
	PgTypeArrCstring:       "PgTypeArrCstring",
	PgTypeTimetz:           "PgTypeTimetz",
	PgTypeArrTimetz:        "PgTypeArrTimetz",
	PgTypeBit:              "PgTypeBit",
	PgTypeArrBit:           "PgTypeArrBit",
	PgTypeVarbit:           "PgTypeVarbit",
	PgTypeArrVarbit:        "PgTypeArrVarbit",
	PgTypeNumeric:          "PgTypeNumeric",
	PgTypeRefcursor:        "PgTypeRefcursor",
	PgTypeArrRefcursor:     "PgTypeArrRefcursor",
	PgTypeRegprocedure:     "PgTypeRegprocedure",
	PgTypeRegoper:          "PgTypeRegoper",
	PgTypeRegoperator:      "PgTypeRegoperator",
	PgTypeRegclass:         "PgTypeRegclass",
	PgTypeRegtype:          "PgTypeRegtype",
	PgTypeArrRegprocedure:  "PgTypeArrRegprocedure",
	PgTypeArrRegoper:       "PgTypeArrRegoper",
	PgTypeArrRegoperator:   "PgTypeArrRegoperator",
	PgTypeArrRegclass:      "PgTypeArrRegclass",
	PgTypeArrRegtype:       "PgTypeArrRegtype",
	PgTypeRecord:           "PgTypeRecord",
	PgTypeCstring:          "PgTypeCstring",
	PgTypeAny:              "PgTypeAny",
	PgTypeAnyarray:         "PgTypeAnyarray",
	PgTypeVoid:             "PgTypeVoid",
	PgTypeTrigger:          "PgTypeTrigger",
	PgTypeLanguageHandler:  "PgTypeLanguageHandler",
	PgTypeInternal:         "PgTypeInternal",
	PgTypeOpaque:           "PgTypeOpaque",
	PgTypeAnyelement:       "PgTypeAnyelement",
	PgTypeArrRecord:        "PgTypeArrRecord",
	PgTypeAnynonarray:      "PgTypeAnynonarray",
	PgTypePgAuthid:         "PgTypePgAuthid",
	PgTypePgAuthMembers:    "PgTypePgAuthMembers",
	PgTypeArrTxidSnapshot:  "PgTypeArrTxidSnapshot",
	PgTypeUuid:             "PgTypeUuid",
	PgTypeArrUuid:          "PgTypeArrUuid",
	PgTypeTxidSnapshot:     "PgTypeTxidSnapshot",
	PgTypeFdwHandler:       "PgTypeFdwHandler",
	PgTypePgLsn:            "PgTypePgLsn",
	PgTypeArrPgLsn:         "PgTypeArrPgLsn",
	PgTypeTsmHandler:       "PgTypeTsmHandler",
	PgTypePgNdistinct:      "PgTypePgNdistinct",
	PgTypePgDependencies:   "PgTypePgDependencies",
	PgTypeAnyenum:          "PgTypeAnyenum",
	PgTypeTsvector:         "PgTypeTsvector",
	PgTypeTsquery:          "PgTypeTsquery",
	PgTypeGtsvector:        "PgTypeGtsvector",
	PgTypeArrTsvector:      "PgTypeArrTsvector",
	PgTypeArrGtsvector:     "PgTypeArrGtsvector",
	PgTypeArrTsquery:       "PgTypeArrTsquery",
	PgTypeRegconfig:        "PgTypeRegconfig",
	PgTypeArrRegconfig:     "PgTypeArrRegconfig",
	PgTypeRegdictionary:    "PgTypeRegdictionary",
	PgTypeArrRegdictionary: "PgTypeArrRegdictionary",
	PgTypeJsonb:            "PgTypeJsonb",
	PgTypeArrJsonb:         "PgTypeArrJsonb",
	PgTypeAnyrange:         "PgTypeAnyrange",
	PgTypeEventTrigger:     "PgTypeEventTrigger",
	PgTypeInt4range:        "PgTypeInt4range",
	PgTypeArrInt4range:     "PgTypeArrInt4range",
	PgTypeNumrange:         "PgTypeNumrange",
	PgTypeArrNumrange:      "PgTypeArrNumrange",
	PgTypeTsrange:          "PgTypeTsrange",
	PgTypeArrTsrange:       "PgTypeArrTsrange",
	PgTypeTstzrange:        "PgTypeTstzrange",
	PgTypeArrTstzrange:     "PgTypeArrTstzrange",
	PgTypeDaterange:        "PgTypeDaterange",
	PgTypeArrDaterange:     "PgTypeArrDaterange",
	PgTypeInt8range:        "PgTypeInt8range",
	PgTypeArrInt8range:     "PgTypeArrInt8range",
	PgTypePgShseclabel:     "PgTypePgShseclabel",
	PgTypeRegnamespace:     "PgTypeRegnamespace",
	PgTypeArrRegnamespace:  "PgTypeArrRegnamespace",
	PgTypeRegrole:          "PgTypeRegrole",
	PgTypeArrRegrole:       "PgTypeArrRegrole",
	PgTypePgSubscription:   "PgTypePgSubscription",
}
