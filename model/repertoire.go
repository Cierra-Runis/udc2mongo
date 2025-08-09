package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// See: https://unicode.org/reports/tr42/#d1e2832
//
// WARNING: The feature of [Group] is not implemented, because the XML does not contain <group>.
//
// However, we can't commit that it will never contained at XML file in the future.
//
// [Group]: https://unicode.org/reports/tr42/#group
type Repertoire struct {
	Reserved     []CodePoint `xml:"reserved" json:"reserved,omitempty"`         // https://unicode.org/reports/tr42/#d1e2899
	Noncharacter []CodePoint `xml:"noncharacter" json:"noncharacter,omitempty"` // https://unicode.org/reports/tr42/#d1e2899
	Surrogate    []CodePoint `xml:"surrogate" json:"surrogate,omitempty"`       // https://unicode.org/reports/tr42/#d1e2899
	CodePoints   []CodePoint `xml:"char" json:"code_points,omitempty"`          // https://unicode.org/reports/tr42/#d1e2899
}

// See: https://unicode.org/reports/tr42/#d1e2899
type CodePoint struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`

	CP      string `xml:"cp,attr" bson:"cp" json:"cp,omitempty"`                   // https://unicode.org/reports/tr42/#d1e2857
	FirstCP string `xml:"first-cp,attr" bson:"first_cp" json:"first_cp,omitempty"` // https://unicode.org/reports/tr42/#d1e2857
	LastCP  string `xml:"last-cp,attr" bson:"last_cp" json:"last_cp,omitempty"`    // https://unicode.org/reports/tr42/#d1e2857

	CodePointProperties `bson:",inline"` // See: https://unicode.org/reports/tr42/lp:d1e2887
}

// See: https://unicode.org/reports/tr42/#d1e3019
type CodePointProperties struct {
	AgeProperties           `bson:",inline"` // See: https://unicode.org/reports/tr42/#d1e3048
	NameProperties          `bson:",inline"` // See: https://unicode.org/reports/tr42/#d1e3071
	NameAliases             []NameAlias      `xml:"name-alias" bson:"name_aliases" json:"name_aliases,omitempty"`      // See: https://unicode.org/reports/tr42/#d1e3145
	Block                   string           `xml:"blk,attr" bson:"block" json:"block,omitempty"`                      // See: https://unicode.org/reports/tr42/#d1e3168
	GeneralCategory         string           `xml:"gc,attr" bson:"general_category" json:"general_category,omitempty"` // See: https://unicode.org/reports/tr42/#d1e3191
	CombiningClass          int              `xml:"ccc,attr" bson:"combining_class" json:"combining_class,omitempty"`  // See: https://unicode.org/reports/tr42/#d1e3215
	BidiProperties          `bson:",inline"` // See: https://unicode.org/reports/tr42/#d1e3241
	DecompositionProperties `bson:",inline"` // See: https://unicode.org/reports/tr42/#d1e3332
	NumericProperties       `bson:",inline"` // See: https://unicode.org/reports/tr42/#d1e3393
	JoiningProperties       `bson:",inline"` // See: https://unicode.org/reports/tr42/#d1e3422
	LineBreak               string           `xml:"lb,attr" bson:"line_break" json:"line_break,omitempty"`             // See: https://unicode.org/reports/tr42/#d1e3467
	EastAsianWidth          string           `xml:"ea,attr" bson:"east_asian_width" json:"east_asian_width,omitempty"` // See: https://unicode.org/reports/tr42/#d1e3491
	CaseProperties          `bson:",inline"` // See: https://unicode.org/reports/tr42/#d1e3514
	ScriptProperties        `bson:",inline"` // See: https://unicode.org/reports/tr42/#d1e3614

	ISOComment string `xml:"isc,attr" bson:"iso_comment" json:"iso_comment,omitempty"`

	// Hangul 属性
	HangulSyllableType string `xml:"hst,attr" bson:"hangul_syllable_type" json:"hangul_syllable_type,omitempty"`
	JamoShortName      string `xml:"JSN,attr" bson:"jamo_short_name" json:"jamo_short_name,omitempty"`

	// Indic 属性
	IndicSyllabicCategory   string `xml:"InSC,attr" bson:"indic_syllabic_category" json:"indic_syllabic_category,omitempty"`
	IndicMatraCategory      string `xml:"InMC,attr" bson:"indic_matra_category" json:"indic_matra_category,omitempty"`
	IndicPositionalCategory string `xml:"InPC,attr" bson:"indic_positional_category" json:"indic_positional_category,omitempty"`
	IndicConjunctBreak      string `xml:"InCB,attr" bson:"indic_conjunct_break" json:"indic_conjunct_break,omitempty"`

	// 标识符属性
	IDStart              UCDBool `xml:"IDS,attr" bson:"id_start" json:"id_start,omitempty"`
	OtherIDStart         UCDBool `xml:"OIDS,attr" bson:"other_id_start" json:"other_id_start,omitempty"`
	XIDStart             UCDBool `xml:"XIDS,attr" bson:"xid_start" json:"xid_start,omitempty"`
	IDContinue           UCDBool `xml:"IDC,attr" bson:"id_continue" json:"id_continue,omitempty"`
	OtherIDContinue      UCDBool `xml:"OIDC,attr" bson:"other_id_continue" json:"other_id_continue,omitempty"`
	XIDContinue          UCDBool `xml:"XIDC,attr" bson:"xid_continue" json:"xid_continue,omitempty"`
	IDCompatMathStart    UCDBool `xml:"ID_Compat_Math_Start,attr" bson:"id_compat_math_start" json:"id_compat_math_start,omitempty"`
	IDCompatMathContinue UCDBool `xml:"ID_Compat_Math_Continue,attr" bson:"id_compat_math_continue" json:"id_compat_math_continue,omitempty"`

	// 模式属性
	PatternSyntax     UCDBool `xml:"Pat_Syn,attr" bson:"pattern_syntax" json:"pattern_syntax,omitempty"`
	PatternWhiteSpace UCDBool `xml:"Pat_WS,attr" bson:"pattern_white_space" json:"pattern_white_space,omitempty"`

	// 标点属性
	Dash                UCDBool `xml:"Dash,attr" bson:"dash" json:"dash,omitempty"`
	Hyphen              UCDBool `xml:"Hyphen,attr" bson:"hyphen" json:"hyphen,omitempty"`
	QuotationMark       UCDBool `xml:"QMark,attr" bson:"quotation_mark" json:"quotation_mark,omitempty"`
	TerminalPunctuation UCDBool `xml:"Term,attr" bson:"terminal_punctuation" json:"terminal_punctuation,omitempty"`
	SentenceTerminal    UCDBool `xml:"STerm,attr" bson:"sentence_terminal" json:"sentence_terminal,omitempty"`

	// 变音符号属性
	Diacritic                  UCDBool `xml:"Dia,attr" bson:"diacritic" json:"diacritic,omitempty"`
	Extender                   UCDBool `xml:"Ext,attr" bson:"extender" json:"extender,omitempty"`
	PrependedConcatenationMark UCDBool `xml:"PCM,attr" bson:"prepended_concatenation_mark" json:"prepended_concatenation_mark,omitempty"`

	// 字符属性
	Alphabetic            UCDBool `xml:"Alpha,attr" bson:"alphabetic" json:"alphabetic,omitempty"`
	OtherAlphabetic       UCDBool `xml:"OAlpha,attr" bson:"other_alphabetic" json:"other_alphabetic,omitempty"`
	Math                  UCDBool `xml:"Math,attr" bson:"math" json:"math,omitempty"`
	OtherMath             UCDBool `xml:"OMath,attr" bson:"other_math" json:"other_math,omitempty"`
	HexDigit              UCDBool `xml:"Hex,attr" bson:"hex_digit" json:"hex_digit,omitempty"`
	ASCIIHexDigit         UCDBool `xml:"AHex,attr" bson:"ascii_hex_digit" json:"ascii_hex_digit,omitempty"`
	DefaultIgnorable      UCDBool `xml:"DI,attr" bson:"default_ignorable" json:"default_ignorable,omitempty"`
	OtherDefaultIgnorable UCDBool `xml:"ODI,attr" bson:"other_default_ignorable" json:"other_default_ignorable,omitempty"`
	LogicalOrderException UCDBool `xml:"LOE,attr" bson:"logical_order_exception" json:"logical_order_exception,omitempty"`
	WhiteSpace            UCDBool `xml:"WSpace,attr" bson:"white_space" json:"white_space,omitempty"`

	// 方向属性
	VerticalOrientation string  `xml:"vo,attr" bson:"vertical_orientation" json:"vertical_orientation,omitempty"`
	RegionalIndicator   UCDBool `xml:"RI,attr" bson:"regional_indicator" json:"regional_indicator,omitempty"`

	// 图形属性
	GraphemeBase        UCDBool `xml:"Gr_Base,attr" bson:"grapheme_base" json:"grapheme_base,omitempty"`
	GraphemeExtend      UCDBool `xml:"Gr_Ext,attr" bson:"grapheme_extend" json:"grapheme_extend,omitempty"`
	OtherGraphemeExtend UCDBool `xml:"OGr_Ext,attr" bson:"other_grapheme_extend" json:"other_grapheme_extend,omitempty"`
	GraphemeLink        UCDBool `xml:"Gr_Link,attr" bson:"grapheme_link" json:"grapheme_link,omitempty"`

	// 断点属性
	GraphemeClusterBreak string `xml:"GCB,attr" bson:"grapheme_cluster_break" json:"grapheme_cluster_break,omitempty"`
	WordBreak            string `xml:"WB,attr" bson:"word_break" json:"word_break,omitempty"`
	SentenceBreak        string `xml:"SB,attr" bson:"sentence_break" json:"sentence_break,omitempty"`

	// 表意字符属性
	Ideographic                UCDBool `xml:"Ideo,attr" bson:"ideographic" json:"ideographic,omitempty"`
	UnifiedIdeograph           UCDBool `xml:"UIdeo,attr" bson:"unified_ideograph" json:"unified_ideograph,omitempty"`
	EquivalentUnifiedIdeograph string  `xml:"EqUIdeo,attr" bson:"equivalent_unified_ideograph" json:"equivalent_unified_ideograph,omitempty"`
	IDSBinaryOperator          UCDBool `xml:"IDSB,attr" bson:"ids_binary_operator" json:"ids_binary_operator,omitempty"`
	IDSTrinaryOperator         UCDBool `xml:"IDST,attr" bson:"ids_trinary_operator" json:"ids_trinary_operator,omitempty"`
	IDSUnaryOperator           UCDBool `xml:"IDSU,attr" bson:"ids_unary_operator" json:"ids_unary_operator,omitempty"`
	Radical                    UCDBool `xml:"Radical,attr" bson:"radical" json:"radical,omitempty"`

	// 其他属性
	Deprecated        UCDBool `xml:"Dep,attr" bson:"deprecated" json:"deprecated,omitempty"`
	VariationSelector UCDBool `xml:"VS,attr" bson:"variation_selector" json:"variation_selector,omitempty"`
	Noncharacter      UCDBool `xml:"NChar,attr" bson:"noncharacter" json:"noncharacter,omitempty"`

	// Emoji 属性
	Emoji                UCDBool `xml:"Emoji,attr" bson:"emoji" json:"emoji,omitempty"`
	EmojiPresentation    UCDBool `xml:"EPres,attr" bson:"emoji_presentation" json:"emoji_presentation,omitempty"`
	EmojiModifier        UCDBool `xml:"EMod,attr" bson:"emoji_modifier" json:"emoji_modifier,omitempty"`
	EmojiModifierBase    UCDBool `xml:"EBase,attr" bson:"emoji_modifier_base" json:"emoji_modifier_base,omitempty"`
	EmojiComponent       UCDBool `xml:"EComp,attr" bson:"emoji_component" json:"emoji_component,omitempty"`
	ExtendedPictographic UCDBool `xml:"ExtPict,attr" bson:"extended_pictographic" json:"extended_pictographic,omitempty"`

	// Unihan 属性 (简化部分重要的)
	KDefinition         string `xml:"kDefinition,attr" bson:"k_definition" json:"k_definition,omitempty"`
	KMandarin           string `xml:"kMandarin,attr" bson:"k_mandarin" json:"k_mandarin,omitempty"`
	KCantonese          string `xml:"kCantonese,attr" bson:"k_cantonese" json:"k_cantonese,omitempty"`
	KJapaneseKun        string `xml:"kJapaneseKun,attr" bson:"k_japanese_kun" json:"k_japanese_kun,omitempty"`
	KJapaneseOn         string `xml:"kJapaneseOn,attr" bson:"k_japanese_on" json:"k_japanese_on,omitempty"`
	KKorean             string `xml:"kKorean,attr" bson:"k_korean" json:"k_korean,omitempty"`
	KVietnamese         string `xml:"kVietnamese,attr" bson:"k_vietnamese" json:"k_vietnamese,omitempty"`
	KTotalStrokes       string `xml:"kTotalStrokes,attr" bson:"k_total_strokes" json:"k_total_strokes,omitempty"`
	KSimplifiedVariant  string `xml:"kSimplifiedVariant,attr" bson:"k_simplified_variant" json:"k_simplified_variant,omitempty"`
	KTraditionalVariant string `xml:"kTraditionalVariant,attr" bson:"k_traditional_variant" json:"k_traditional_variant,omitempty"`

	// 时间戳
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// See: https://unicode.org/reports/tr42/#d1e3048
type AgeProperties struct {
	Age string `xml:"age,attr" bson:"age" json:"age,omitempty"`
}

// See: https://unicode.org/reports/tr42/#d1e3071
type NameProperties struct {
	Name  string `xml:"na,attr" bson:"name" json:"name,omitempty"`
	Name1 string `xml:"na1,attr" bson:"name1" json:"name1,omitempty"`
}

// See: https://unicode.org/reports/tr42/#d1e3145
type NameAlias struct {
	Alias string `xml:"alias,attr" bson:"alias" json:"alias,omitempty"`
	Type  string `xml:"type,attr" bson:"type" json:"type,omitempty"`
}

// See: https://unicode.org/reports/tr42/#d1e3241
type BidiProperties struct {
	BidiClass             string  `xml:"bc,attr" bson:"bidi_class" json:"bidi_class,omitempty"`                              // See: https://unicode.org/reports/tr42/#lp:d1e3148
	BidiMirrored          UCDBool `xml:"Bidi_M,attr" bson:"bidi_mirrored" json:"bidi_mirrored,omitempty"`                    // See: https://unicode.org/reports/tr42/#lp:d1e3157
	BidiMirroringGlyph    string  `xml:"bmg,attr" bson:"bidi_mirroring_glyph" json:"bidi_mirroring_glyph,omitempty"`         // See: https://unicode.org/reports/tr42/#lp:d1e3167
	BidiControl           UCDBool `xml:"Bidi_C,attr" bson:"bidi_control" json:"bidi_control,omitempty"`                      // See: https://unicode.org/reports/tr42/#lp:d1e3179
	BidiPairedBracketType string  `xml:"bpt,attr" bson:"bidi_paired_bracket_type" json:"bidi_paired_bracket_type,omitempty"` // See: https://unicode.org/reports/tr42/#lp:d1e3192
	BidiPairedBracket     string  `xml:"bpb,attr" bson:"bidi_paired_bracket" json:"bidi_paired_bracket,omitempty"`           // See: https://unicode.org/reports/tr42/#lp:d1e3192
}

// See: https://unicode.org/reports/tr42/#d1e3332
type DecompositionProperties struct {
	DecompositionType    string `xml:"dt,attr" bson:"decomposition_type" json:"decomposition_type,omitempty"`       // See: https://unicode.org/reports/tr42/#lp:d1e3215
	DecompositionMapping string `xml:"dm,attr" bson:"decomposition_mapping" json:"decomposition_mapping,omitempty"` // See: https://unicode.org/reports/tr42/#lp:d1e3215

	CompositionExclusion     UCDBool `xml:"CE,attr" bson:"composition_exclusion" json:"composition_exclusion,omitempty"`                // See: https://unicode.org/reports/tr42/#lp:d1e3228
	FullCompositionExclusion UCDBool `xml:"Comp_Ex,attr" bson:"full_composition_exclusion" json:"full_composition_exclusion,omitempty"` // See: https://unicode.org/reports/tr42/#lp:d1e3228

	NFC_QC  string  `xml:"NFC_QC,attr" bson:"nfc_qc" json:"nfc_qc,omitempty"`    // See: https://unicode.org/reports/tr42/#lp:d1e3234
	NFD_QC  string  `xml:"NFD_QC,attr" bson:"nfd_qc" json:"nfd_qc,omitempty"`    // See: https://unicode.org/reports/tr42/#lp:d1e3234
	NFKC_QC string  `xml:"NFKC_QC,attr" bson:"nfkc_qc" json:"nfkc_qc,omitempty"` // See: https://unicode.org/reports/tr42/#lp:d1e3234
	NFKD_QC string  `xml:"NFKD_QC,attr" bson:"nfkd_qc" json:"nfkd_qc,omitempty"` // See: https://unicode.org/reports/tr42/#lp:d1e3234
	XO_NFC  UCDBool `xml:"XO_NFC,attr" bson:"xo_nfc" json:"xo_nfc,omitempty"`    // See: https://unicode.org/reports/tr42/#lp:d1e3234
	XO_NFD  UCDBool `xml:"XO_NFD,attr" bson:"xo_nfd" json:"xo_nfd,omitempty"`    // See: https://unicode.org/reports/tr42/#lp:d1e3234
	XO_NFKC UCDBool `xml:"XO_NFKC,attr" bson:"xo_nfkc" json:"xo_nfkc,omitempty"` // See: https://unicode.org/reports/tr42/#lp:d1e3234
	XO_NFKD UCDBool `xml:"XO_NFKD,attr" bson:"xo_nfkd" json:"xo_nfkd,omitempty"` // See: https://unicode.org/reports/tr42/#lp:d1e3234
	FC_NFKC string  `xml:"FC_NFKC,attr" bson:"fc_nfkc" json:"fc_nfkc,omitempty"` // See: https://unicode.org/reports/tr42/#lp:d1e3234
}

// See: https://unicode.org/reports/tr42/#d1e3393
type NumericProperties struct {
	// See: https://unicode.org/reports/tr42/#lp:d1e3258
	NumericType string `xml:"nt,attr" bson:"numeric_type" json:"numeric_type,omitempty"`
	// See: https://unicode.org/reports/tr42/#lp:d1e3258
	// TODO: Check why it's list
	// 	attribute nv { "NaN" | list { xsd:string { pattern = "-?[0-9]+(/[0-9]+)?" } +}}?
	NumericValue string `xml:"nv,attr" bson:"numeric_value" json:"numeric_value,omitempty"`
}

// See: https://unicode.org/reports/tr42/#d1e3422
type JoiningProperties struct {
	JoiningType  string  `xml:"jt,attr" bson:"joining_type" json:"joining_type,omitempty"`     // See: https://unicode.org/reports/tr42/#lp:d1e3281
	JoiningGroup string  `xml:"jg,attr" bson:"joining_group" json:"joining_group,omitempty"`   // See: https://unicode.org/reports/tr42/#lp:d1e3281
	JoinControl  UCDBool `xml:"Join_C,attr" bson:"join_control" json:"join_control,omitempty"` // See: https://unicode.org/reports/tr42/#lp:d1e3291
}

// See: https://unicode.org/reports/tr42/#d1e3514
type CaseProperties struct {
	Uppercase      UCDBool `xml:"Upper,attr" bson:"uppercase" json:"uppercase,omitempty"`              // See: https://unicode.org/reports/tr42/#lp:d1e3340
	Lowercase      UCDBool `xml:"Lower,attr" bson:"lowercase" json:"lowercase,omitempty"`              // See: https://unicode.org/reports/tr42/#lp:d1e3340
	OtherUppercase UCDBool `xml:"OUpper,attr" bson:"other_uppercase" json:"other_uppercase,omitempty"` // See: https://unicode.org/reports/tr42/#lp:d1e3340
	OtherLowercase UCDBool `xml:"OLower,attr" bson:"other_lowercase" json:"other_lowercase,omitempty"` // See: https://unicode.org/reports/tr42/#lp:d1e3340

	SimpleUppercase string `xml:"suc,attr" bson:"simple_uppercase" json:"simple_uppercase,omitempty"` // See: https://unicode.org/reports/tr42/#lp:d1e3359
	SimpleLowercase string `xml:"slc,attr" bson:"simple_lowercase" json:"simple_lowercase,omitempty"` // See: https://unicode.org/reports/tr42/#lp:d1e3359
	SimpleTitlecase string `xml:"stc,attr" bson:"simple_titlecase" json:"simple_titlecase,omitempty"` // See: https://unicode.org/reports/tr42/#lp:d1e3359

	UppercaseMapping string `xml:"uc,attr" bson:"uppercase_mapping" json:"uppercase_mapping,omitempty"` // See: https://unicode.org/reports/tr42/#lp:d1e3374
	LowercaseMapping string `xml:"lc,attr" bson:"lowercase_mapping" json:"lowercase_mapping,omitempty"` // See: https://unicode.org/reports/tr42/#lp:d1e3374
	TitlecaseMapping string `xml:"tc,attr" bson:"titlecase_mapping" json:"titlecase_mapping,omitempty"` // See: https://unicode.org/reports/tr42/#lp:d1e3374

	SimpleCaseFolding string `xml:"scf,attr" bson:"simple_case_folding" json:"simple_case_folding,omitempty"` // See: https://unicode.org/reports/tr42/#lp:d1e3386
	CaseFolding       string `xml:"cf,attr" bson:"case_folding" json:"case_folding,omitempty"`                // See: https://unicode.org/reports/tr42/#lp:d1e3386

	CaseIgnorable             UCDBool `xml:"CI,attr" bson:"case_ignorable" json:"case_ignorable,omitempty"`                                // See: https://unicode.org/reports/tr42/#lp:d1e3393
	Cased                     UCDBool `xml:"Cased,attr" bson:"cased" json:"cased,omitempty"`                                               // See: https://unicode.org/reports/tr42/#lp:d1e3393
	ChangesWhenCasefolded     UCDBool `xml:"CWCF,attr" bson:"changes_when_casefolded" json:"changes_when_casefolded,omitempty"`            // See: https://unicode.org/reports/tr42/#lp:d1e3393
	ChangesWhenCasemapped     UCDBool `xml:"CWCM,attr" bson:"changes_when_casemapped" json:"changes_when_casemapped,omitempty"`            // See: https://unicode.org/reports/tr42/#lp:d1e3393
	ChangesWhenLowercased     UCDBool `xml:"CWL,attr" bson:"changes_when_lowercased" json:"changes_when_lowercased,omitempty"`             // See: https://unicode.org/reports/tr42/#lp:d1e3393
	ChangesWhenNFKCCasefolded UCDBool `xml:"CWKCF,attr" bson:"changes_when_nfkc_casefolded" json:"changes_when_nfkc_casefolded,omitempty"` // See: https://unicode.org/reports/tr42/#lp:d1e3393
	ChangesWhenTitlecased     UCDBool `xml:"CWT,attr" bson:"changes_when_titlecased" json:"changes_when_titlecased,omitempty"`             // See: https://unicode.org/reports/tr42/#lp:d1e3393
	ChangesWhenUppercased     UCDBool `xml:"CWU,attr" bson:"changes_when_uppercased" json:"changes_when_uppercased,omitempty"`             // See: https://unicode.org/reports/tr42/#lp:d1e3393
	NFKC_CF                   string  `xml:"NFKC_CF,attr" bson:"nfkc_cf" json:"nfkc_cf,omitempty"`                                         // See: https://unicode.org/reports/tr42/#lp:d1e3393
	NFKC_SCF                  string  `xml:"NFKC_SCF,attr" bson:"nfkc_scf" json:"nfkc_scf,omitempty"`                                      // See: https://unicode.org/reports/tr42/#lp:d1e3393
}

// See: https://unicode.org/reports/tr42/#d1e3614
type ScriptProperties struct {
	Script           string `xml:"sc,attr" bson:"script" json:"script,omitempty"`
	ScriptExtensions string `xml:"scx,attr" bson:"script_extensions" json:"script_extensions,omitempty"`
}
