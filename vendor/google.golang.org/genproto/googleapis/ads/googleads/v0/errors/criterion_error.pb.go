// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v0/errors/criterion_error.proto

package errors // import "google.golang.org/genproto/googleapis/ads/googleads/v0/errors"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// Enum describing possible criterion errors.
type CriterionErrorEnum_CriterionError int32

const (
	// Enum unspecified.
	CriterionErrorEnum_UNSPECIFIED CriterionErrorEnum_CriterionError = 0
	// The received error code is not known in this version.
	CriterionErrorEnum_UNKNOWN CriterionErrorEnum_CriterionError = 1
	// Concrete type of criterion is required for CREATE and UPDATE operations.
	CriterionErrorEnum_CONCRETE_TYPE_REQUIRED CriterionErrorEnum_CriterionError = 2
	// The category requested for exclusion is invalid.
	CriterionErrorEnum_INVALID_EXCLUDED_CATEGORY CriterionErrorEnum_CriterionError = 3
	// Invalid keyword criteria text.
	CriterionErrorEnum_INVALID_KEYWORD_TEXT CriterionErrorEnum_CriterionError = 4
	// Keyword text should be less than 80 chars.
	CriterionErrorEnum_KEYWORD_TEXT_TOO_LONG CriterionErrorEnum_CriterionError = 5
	// Keyword text has too many words.
	CriterionErrorEnum_KEYWORD_HAS_TOO_MANY_WORDS CriterionErrorEnum_CriterionError = 6
	// Keyword text has invalid characters or symbols.
	CriterionErrorEnum_KEYWORD_HAS_INVALID_CHARS CriterionErrorEnum_CriterionError = 7
	// Invalid placement URL.
	CriterionErrorEnum_INVALID_PLACEMENT_URL CriterionErrorEnum_CriterionError = 8
	// Invalid user list criterion.
	CriterionErrorEnum_INVALID_USER_LIST CriterionErrorEnum_CriterionError = 9
	// Invalid user interest criterion.
	CriterionErrorEnum_INVALID_USER_INTEREST CriterionErrorEnum_CriterionError = 10
	// Placement URL has wrong format.
	CriterionErrorEnum_INVALID_FORMAT_FOR_PLACEMENT_URL CriterionErrorEnum_CriterionError = 11
	// Placement URL is too long.
	CriterionErrorEnum_PLACEMENT_URL_IS_TOO_LONG CriterionErrorEnum_CriterionError = 12
	// Indicates the URL contains an illegal character.
	CriterionErrorEnum_PLACEMENT_URL_HAS_ILLEGAL_CHAR CriterionErrorEnum_CriterionError = 13
	// Indicates the URL contains multiple comma separated URLs.
	CriterionErrorEnum_PLACEMENT_URL_HAS_MULTIPLE_SITES_IN_LINE CriterionErrorEnum_CriterionError = 14
	// Indicates the domain is blacklisted.
	CriterionErrorEnum_PLACEMENT_IS_NOT_AVAILABLE_FOR_TARGETING_OR_EXCLUSION CriterionErrorEnum_CriterionError = 15
	// Invalid topic path.
	CriterionErrorEnum_INVALID_TOPIC_PATH CriterionErrorEnum_CriterionError = 16
	// The YouTube Channel Id is invalid.
	CriterionErrorEnum_INVALID_YOUTUBE_CHANNEL_ID CriterionErrorEnum_CriterionError = 17
	// The YouTube Video Id is invalid.
	CriterionErrorEnum_INVALID_YOUTUBE_VIDEO_ID CriterionErrorEnum_CriterionError = 18
	// Indicates the placement is a YouTube vertical channel, which is no longer
	// supported.
	CriterionErrorEnum_YOUTUBE_VERTICAL_CHANNEL_DEPRECATED CriterionErrorEnum_CriterionError = 19
	// Indicates the placement is a YouTube demographic channel, which is no
	// longer supported.
	CriterionErrorEnum_YOUTUBE_DEMOGRAPHIC_CHANNEL_DEPRECATED CriterionErrorEnum_CriterionError = 20
	// YouTube urls are not supported in Placement criterion. Use YouTubeChannel
	// and YouTubeVideo criterion instead.
	CriterionErrorEnum_YOUTUBE_URL_UNSUPPORTED CriterionErrorEnum_CriterionError = 21
	// Criteria type can not be excluded by the customer, like AOL account type
	// cannot target site type criteria.
	CriterionErrorEnum_CANNOT_EXCLUDE_CRITERIA_TYPE CriterionErrorEnum_CriterionError = 22
	// Criteria type can not be targeted.
	CriterionErrorEnum_CANNOT_ADD_CRITERIA_TYPE CriterionErrorEnum_CriterionError = 23
	// Product filter in the product criteria has invalid characters. Operand
	// and the argument in the filter can not have "==" or "&+".
	CriterionErrorEnum_INVALID_PRODUCT_FILTER CriterionErrorEnum_CriterionError = 24
	// Product filter in the product criteria is translated to a string as
	// operand1==argument1&+operand2==argument2, maximum allowed length for the
	// string is 255 chars.
	CriterionErrorEnum_PRODUCT_FILTER_TOO_LONG CriterionErrorEnum_CriterionError = 25
	// Not allowed to exclude similar user list.
	CriterionErrorEnum_CANNOT_EXCLUDE_SIMILAR_USER_LIST CriterionErrorEnum_CriterionError = 26
	// Not allowed to target a closed user list.
	CriterionErrorEnum_CANNOT_ADD_CLOSED_USER_LIST CriterionErrorEnum_CriterionError = 27
	// Not allowed to add display only UserLists to search only campaigns.
	CriterionErrorEnum_CANNOT_ADD_DISPLAY_ONLY_LISTS_TO_SEARCH_ONLY_CAMPAIGNS CriterionErrorEnum_CriterionError = 28
	// Not allowed to add display only UserLists to search plus campaigns.
	CriterionErrorEnum_CANNOT_ADD_DISPLAY_ONLY_LISTS_TO_SEARCH_CAMPAIGNS CriterionErrorEnum_CriterionError = 29
	// Not allowed to add display only UserLists to shopping campaigns.
	CriterionErrorEnum_CANNOT_ADD_DISPLAY_ONLY_LISTS_TO_SHOPPING_CAMPAIGNS CriterionErrorEnum_CriterionError = 30
	// Not allowed to add User interests to search only campaigns.
	CriterionErrorEnum_CANNOT_ADD_USER_INTERESTS_TO_SEARCH_CAMPAIGNS CriterionErrorEnum_CriterionError = 31
	// Not allowed to set bids for this criterion type in search campaigns
	CriterionErrorEnum_CANNOT_SET_BIDS_ON_CRITERION_TYPE_IN_SEARCH_CAMPAIGNS CriterionErrorEnum_CriterionError = 32
	// Final URLs, URL Templates and CustomParameters cannot be set for the
	// criterion types of Gender, AgeRange, UserList, Placement, MobileApp, and
	// MobileAppCategory in search campaigns and shopping campaigns.
	CriterionErrorEnum_CANNOT_ADD_URLS_TO_CRITERION_TYPE_FOR_CAMPAIGN_TYPE CriterionErrorEnum_CriterionError = 33
	// IP address is not valid.
	CriterionErrorEnum_INVALID_IP_ADDRESS CriterionErrorEnum_CriterionError = 34
	// IP format is not valid.
	CriterionErrorEnum_INVALID_IP_FORMAT CriterionErrorEnum_CriterionError = 35
	// Mobile application is not valid.
	CriterionErrorEnum_INVALID_MOBILE_APP CriterionErrorEnum_CriterionError = 36
	// Mobile application category is not valid.
	CriterionErrorEnum_INVALID_MOBILE_APP_CATEGORY CriterionErrorEnum_CriterionError = 37
	// The CriterionId does not exist or is of the incorrect type.
	CriterionErrorEnum_INVALID_CRITERION_ID CriterionErrorEnum_CriterionError = 38
	// The Criterion is not allowed to be targeted.
	CriterionErrorEnum_CANNOT_TARGET_CRITERION CriterionErrorEnum_CriterionError = 39
	// The criterion is not allowed to be targeted as it is deprecated.
	CriterionErrorEnum_CANNOT_TARGET_OBSOLETE_CRITERION CriterionErrorEnum_CriterionError = 40
	// The CriterionId is not valid for the type.
	CriterionErrorEnum_CRITERION_ID_AND_TYPE_MISMATCH CriterionErrorEnum_CriterionError = 41
	// Distance for the radius for the proximity criterion is invalid.
	CriterionErrorEnum_INVALID_PROXIMITY_RADIUS CriterionErrorEnum_CriterionError = 42
	// Units for the distance for the radius for the proximity criterion is
	// invalid.
	CriterionErrorEnum_INVALID_PROXIMITY_RADIUS_UNITS CriterionErrorEnum_CriterionError = 43
	// Street address in the address is not valid.
	CriterionErrorEnum_INVALID_STREETADDRESS_LENGTH CriterionErrorEnum_CriterionError = 44
	// City name in the address is not valid.
	CriterionErrorEnum_INVALID_CITYNAME_LENGTH CriterionErrorEnum_CriterionError = 45
	// Region code in the address is not valid.
	CriterionErrorEnum_INVALID_REGIONCODE_LENGTH CriterionErrorEnum_CriterionError = 46
	// Region name in the address is not valid.
	CriterionErrorEnum_INVALID_REGIONNAME_LENGTH CriterionErrorEnum_CriterionError = 47
	// Postal code in the address is not valid.
	CriterionErrorEnum_INVALID_POSTALCODE_LENGTH CriterionErrorEnum_CriterionError = 48
	// Country code in the address is not valid.
	CriterionErrorEnum_INVALID_COUNTRY_CODE CriterionErrorEnum_CriterionError = 49
	// Latitude for the GeoPoint is not valid.
	CriterionErrorEnum_INVALID_LATITUDE CriterionErrorEnum_CriterionError = 50
	// Longitude for the GeoPoint is not valid.
	CriterionErrorEnum_INVALID_LONGITUDE CriterionErrorEnum_CriterionError = 51
	// The Proximity input is not valid. Both address and geoPoint cannot be
	// null.
	CriterionErrorEnum_PROXIMITY_GEOPOINT_AND_ADDRESS_BOTH_CANNOT_BE_NULL CriterionErrorEnum_CriterionError = 52
	// The Proximity address cannot be geocoded to a valid lat/long.
	CriterionErrorEnum_INVALID_PROXIMITY_ADDRESS CriterionErrorEnum_CriterionError = 53
	// User domain name is not valid.
	CriterionErrorEnum_INVALID_USER_DOMAIN_NAME CriterionErrorEnum_CriterionError = 54
	// Length of serialized criterion parameter exceeded size limit.
	CriterionErrorEnum_CRITERION_PARAMETER_TOO_LONG CriterionErrorEnum_CriterionError = 55
	// Time interval in the AdSchedule overlaps with another AdSchedule.
	CriterionErrorEnum_AD_SCHEDULE_TIME_INTERVALS_OVERLAP CriterionErrorEnum_CriterionError = 56
	// AdSchedule time interval cannot span multiple days.
	CriterionErrorEnum_AD_SCHEDULE_INTERVAL_CANNOT_SPAN_MULTIPLE_DAYS CriterionErrorEnum_CriterionError = 57
	// AdSchedule time interval specified is invalid, endTime cannot be earlier
	// than startTime.
	CriterionErrorEnum_AD_SCHEDULE_INVALID_TIME_INTERVAL CriterionErrorEnum_CriterionError = 58
	// The number of AdSchedule entries in a day exceeds the limit.
	CriterionErrorEnum_AD_SCHEDULE_EXCEEDED_INTERVALS_PER_DAY_LIMIT CriterionErrorEnum_CriterionError = 59
	// CriteriaId does not match the interval of the AdSchedule specified.
	CriterionErrorEnum_AD_SCHEDULE_CRITERION_ID_MISMATCHING_FIELDS CriterionErrorEnum_CriterionError = 60
	// Cannot set bid modifier for this criterion type.
	CriterionErrorEnum_CANNOT_BID_MODIFY_CRITERION_TYPE CriterionErrorEnum_CriterionError = 61
	// Cannot bid modify criterion, since it is opted out of the campaign.
	CriterionErrorEnum_CANNOT_BID_MODIFY_CRITERION_CAMPAIGN_OPTED_OUT CriterionErrorEnum_CriterionError = 62
	// Cannot set bid modifier for a negative criterion.
	CriterionErrorEnum_CANNOT_BID_MODIFY_NEGATIVE_CRITERION CriterionErrorEnum_CriterionError = 63
	// Bid Modifier already exists. Use SET operation to update.
	CriterionErrorEnum_BID_MODIFIER_ALREADY_EXISTS CriterionErrorEnum_CriterionError = 64
	// Feed Id is not allowed in these Location Groups.
	CriterionErrorEnum_FEED_ID_NOT_ALLOWED CriterionErrorEnum_CriterionError = 65
	// The account may not use the requested criteria type. For example, some
	// accounts are restricted to keywords only.
	CriterionErrorEnum_ACCOUNT_INELIGIBLE_FOR_CRITERIA_TYPE CriterionErrorEnum_CriterionError = 66
	// The requested criteria type cannot be used with campaign or ad group
	// bidding strategy.
	CriterionErrorEnum_CRITERIA_TYPE_INVALID_FOR_BIDDING_STRATEGY CriterionErrorEnum_CriterionError = 67
	// The Criterion is not allowed to be excluded.
	CriterionErrorEnum_CANNOT_EXCLUDE_CRITERION CriterionErrorEnum_CriterionError = 68
	// The criterion is not allowed to be removed. For example, we cannot remove
	// any of the device criterion.
	CriterionErrorEnum_CANNOT_REMOVE_CRITERION CriterionErrorEnum_CriterionError = 69
	// The combined length of product dimension values of the product scope
	// criterion is too long.
	CriterionErrorEnum_PRODUCT_SCOPE_TOO_LONG CriterionErrorEnum_CriterionError = 70
	// Product scope contains too many dimensions.
	CriterionErrorEnum_PRODUCT_SCOPE_TOO_MANY_DIMENSIONS CriterionErrorEnum_CriterionError = 71
	// The combined length of product dimension values of the product partition
	// criterion is too long.
	CriterionErrorEnum_PRODUCT_PARTITION_TOO_LONG CriterionErrorEnum_CriterionError = 72
	// Product partition contains too many dimensions.
	CriterionErrorEnum_PRODUCT_PARTITION_TOO_MANY_DIMENSIONS CriterionErrorEnum_CriterionError = 73
	// The product dimension is invalid (e.g. dimension contains illegal value,
	// dimension type is represented with wrong class, etc). Product dimension
	// value can not contain "==" or "&+".
	CriterionErrorEnum_INVALID_PRODUCT_DIMENSION CriterionErrorEnum_CriterionError = 74
	// Product dimension type is either invalid for campaigns of this type or
	// cannot be used in the current context. BIDDING_CATEGORY_Lx and
	// PRODUCT_TYPE_Lx product dimensions must be used in ascending order of
	// their levels: L1, L2, L3, L4, L5... The levels must be specified
	// sequentially and start from L1. Furthermore, an "others" product
	// partition cannot be subdivided with a dimension of the same type but of a
	// higher level ("others" BIDDING_CATEGORY_L3 can be subdivided with BRAND
	// but not with BIDDING_CATEGORY_L4).
	CriterionErrorEnum_INVALID_PRODUCT_DIMENSION_TYPE CriterionErrorEnum_CriterionError = 75
	// Bidding categories do not form a valid path in the Shopping bidding
	// category taxonomy.
	CriterionErrorEnum_INVALID_PRODUCT_BIDDING_CATEGORY CriterionErrorEnum_CriterionError = 76
	// ShoppingSetting must be added to the campaign before ProductScope
	// criteria can be added.
	CriterionErrorEnum_MISSING_SHOPPING_SETTING CriterionErrorEnum_CriterionError = 77
	// Matching function is invalid.
	CriterionErrorEnum_INVALID_MATCHING_FUNCTION CriterionErrorEnum_CriterionError = 78
	// Filter parameters not allowed for location groups targeting.
	CriterionErrorEnum_LOCATION_FILTER_NOT_ALLOWED CriterionErrorEnum_CriterionError = 79
	// Given location filter parameter is invalid for location groups targeting.
	CriterionErrorEnum_LOCATION_FILTER_INVALID CriterionErrorEnum_CriterionError = 80
	// Criteria type cannot be associated with a campaign and its ad group(s)
	// simultaneously.
	CriterionErrorEnum_CANNOT_ATTACH_CRITERIA_AT_CAMPAIGN_AND_ADGROUP CriterionErrorEnum_CriterionError = 81
	// Range represented by hotel length of stay's min nights and max nights
	// overlaps with an existing criterion.
	CriterionErrorEnum_HOTEL_LENGTH_OF_STAY_OVERLAPS_WITH_EXISTING_CRITERION CriterionErrorEnum_CriterionError = 82
	// Range represented by hotel advance booking window's min days and max days
	// overlaps with an existing criterion.
	CriterionErrorEnum_HOTEL_ADVANCE_BOOKING_WINDOW_OVERLAPS_WITH_EXISTING_CRITERION CriterionErrorEnum_CriterionError = 83
	// The field is not allowed to be set when the negative field is set to
	// true, e.g. we don't allow bids in negative ad group or campaign criteria.
	CriterionErrorEnum_FIELD_INCOMPATIBLE_WITH_NEGATIVE_TARGETING CriterionErrorEnum_CriterionError = 84
	// The combination of operand and operator in webpage condition is invalid.
	CriterionErrorEnum_INVALID_WEBPAGE_CONDITION CriterionErrorEnum_CriterionError = 85
	// The URL of webpage condition is invalid.
	CriterionErrorEnum_INVALID_WEBPAGE_CONDITION_URL CriterionErrorEnum_CriterionError = 86
	// The URL of webpage condition cannot be empty or contain white space.
	CriterionErrorEnum_WEBPAGE_CONDITION_URL_CANNOT_BE_EMPTY CriterionErrorEnum_CriterionError = 87
	// The URL of webpage condition contains an unsupported protocol.
	CriterionErrorEnum_WEBPAGE_CONDITION_URL_UNSUPPORTED_PROTOCOL CriterionErrorEnum_CriterionError = 88
	// The URL of webpage condition cannot be an IP address.
	CriterionErrorEnum_WEBPAGE_CONDITION_URL_CANNOT_BE_IP_ADDRESS CriterionErrorEnum_CriterionError = 89
	// The domain of the URL is not consistent with the domain in campaign
	// setting.
	CriterionErrorEnum_WEBPAGE_CONDITION_URL_DOMAIN_NOT_CONSISTENT_WITH_CAMPAIGN_SETTING CriterionErrorEnum_CriterionError = 90
	// The URL of webpage condition cannot be a public suffix itself.
	CriterionErrorEnum_WEBPAGE_CONDITION_URL_CANNOT_BE_PUBLIC_SUFFIX CriterionErrorEnum_CriterionError = 91
	// The URL of webpage condition has an invalid public suffix.
	CriterionErrorEnum_WEBPAGE_CONDITION_URL_INVALID_PUBLIC_SUFFIX CriterionErrorEnum_CriterionError = 92
	// Value track parameter is not supported in webpage condition URL.
	CriterionErrorEnum_WEBPAGE_CONDITION_URL_VALUE_TRACK_VALUE_NOT_SUPPORTED CriterionErrorEnum_CriterionError = 93
	// Only one URL-EQUALS webpage condition is allowed in a webpage
	// criterion and it cannot be combined with other conditions.
	CriterionErrorEnum_WEBPAGE_CRITERION_URL_EQUALS_CAN_HAVE_ONLY_ONE_CONDITION CriterionErrorEnum_CriterionError = 94
	// A webpage criterion cannot be added to a non-DSA ad group.
	CriterionErrorEnum_WEBPAGE_CRITERION_NOT_SUPPORTED_ON_NON_DSA_AD_GROUP CriterionErrorEnum_CriterionError = 95
)

var CriterionErrorEnum_CriterionError_name = map[int32]string{
	0:  "UNSPECIFIED",
	1:  "UNKNOWN",
	2:  "CONCRETE_TYPE_REQUIRED",
	3:  "INVALID_EXCLUDED_CATEGORY",
	4:  "INVALID_KEYWORD_TEXT",
	5:  "KEYWORD_TEXT_TOO_LONG",
	6:  "KEYWORD_HAS_TOO_MANY_WORDS",
	7:  "KEYWORD_HAS_INVALID_CHARS",
	8:  "INVALID_PLACEMENT_URL",
	9:  "INVALID_USER_LIST",
	10: "INVALID_USER_INTEREST",
	11: "INVALID_FORMAT_FOR_PLACEMENT_URL",
	12: "PLACEMENT_URL_IS_TOO_LONG",
	13: "PLACEMENT_URL_HAS_ILLEGAL_CHAR",
	14: "PLACEMENT_URL_HAS_MULTIPLE_SITES_IN_LINE",
	15: "PLACEMENT_IS_NOT_AVAILABLE_FOR_TARGETING_OR_EXCLUSION",
	16: "INVALID_TOPIC_PATH",
	17: "INVALID_YOUTUBE_CHANNEL_ID",
	18: "INVALID_YOUTUBE_VIDEO_ID",
	19: "YOUTUBE_VERTICAL_CHANNEL_DEPRECATED",
	20: "YOUTUBE_DEMOGRAPHIC_CHANNEL_DEPRECATED",
	21: "YOUTUBE_URL_UNSUPPORTED",
	22: "CANNOT_EXCLUDE_CRITERIA_TYPE",
	23: "CANNOT_ADD_CRITERIA_TYPE",
	24: "INVALID_PRODUCT_FILTER",
	25: "PRODUCT_FILTER_TOO_LONG",
	26: "CANNOT_EXCLUDE_SIMILAR_USER_LIST",
	27: "CANNOT_ADD_CLOSED_USER_LIST",
	28: "CANNOT_ADD_DISPLAY_ONLY_LISTS_TO_SEARCH_ONLY_CAMPAIGNS",
	29: "CANNOT_ADD_DISPLAY_ONLY_LISTS_TO_SEARCH_CAMPAIGNS",
	30: "CANNOT_ADD_DISPLAY_ONLY_LISTS_TO_SHOPPING_CAMPAIGNS",
	31: "CANNOT_ADD_USER_INTERESTS_TO_SEARCH_CAMPAIGNS",
	32: "CANNOT_SET_BIDS_ON_CRITERION_TYPE_IN_SEARCH_CAMPAIGNS",
	33: "CANNOT_ADD_URLS_TO_CRITERION_TYPE_FOR_CAMPAIGN_TYPE",
	34: "INVALID_IP_ADDRESS",
	35: "INVALID_IP_FORMAT",
	36: "INVALID_MOBILE_APP",
	37: "INVALID_MOBILE_APP_CATEGORY",
	38: "INVALID_CRITERION_ID",
	39: "CANNOT_TARGET_CRITERION",
	40: "CANNOT_TARGET_OBSOLETE_CRITERION",
	41: "CRITERION_ID_AND_TYPE_MISMATCH",
	42: "INVALID_PROXIMITY_RADIUS",
	43: "INVALID_PROXIMITY_RADIUS_UNITS",
	44: "INVALID_STREETADDRESS_LENGTH",
	45: "INVALID_CITYNAME_LENGTH",
	46: "INVALID_REGIONCODE_LENGTH",
	47: "INVALID_REGIONNAME_LENGTH",
	48: "INVALID_POSTALCODE_LENGTH",
	49: "INVALID_COUNTRY_CODE",
	50: "INVALID_LATITUDE",
	51: "INVALID_LONGITUDE",
	52: "PROXIMITY_GEOPOINT_AND_ADDRESS_BOTH_CANNOT_BE_NULL",
	53: "INVALID_PROXIMITY_ADDRESS",
	54: "INVALID_USER_DOMAIN_NAME",
	55: "CRITERION_PARAMETER_TOO_LONG",
	56: "AD_SCHEDULE_TIME_INTERVALS_OVERLAP",
	57: "AD_SCHEDULE_INTERVAL_CANNOT_SPAN_MULTIPLE_DAYS",
	58: "AD_SCHEDULE_INVALID_TIME_INTERVAL",
	59: "AD_SCHEDULE_EXCEEDED_INTERVALS_PER_DAY_LIMIT",
	60: "AD_SCHEDULE_CRITERION_ID_MISMATCHING_FIELDS",
	61: "CANNOT_BID_MODIFY_CRITERION_TYPE",
	62: "CANNOT_BID_MODIFY_CRITERION_CAMPAIGN_OPTED_OUT",
	63: "CANNOT_BID_MODIFY_NEGATIVE_CRITERION",
	64: "BID_MODIFIER_ALREADY_EXISTS",
	65: "FEED_ID_NOT_ALLOWED",
	66: "ACCOUNT_INELIGIBLE_FOR_CRITERIA_TYPE",
	67: "CRITERIA_TYPE_INVALID_FOR_BIDDING_STRATEGY",
	68: "CANNOT_EXCLUDE_CRITERION",
	69: "CANNOT_REMOVE_CRITERION",
	70: "PRODUCT_SCOPE_TOO_LONG",
	71: "PRODUCT_SCOPE_TOO_MANY_DIMENSIONS",
	72: "PRODUCT_PARTITION_TOO_LONG",
	73: "PRODUCT_PARTITION_TOO_MANY_DIMENSIONS",
	74: "INVALID_PRODUCT_DIMENSION",
	75: "INVALID_PRODUCT_DIMENSION_TYPE",
	76: "INVALID_PRODUCT_BIDDING_CATEGORY",
	77: "MISSING_SHOPPING_SETTING",
	78: "INVALID_MATCHING_FUNCTION",
	79: "LOCATION_FILTER_NOT_ALLOWED",
	80: "LOCATION_FILTER_INVALID",
	81: "CANNOT_ATTACH_CRITERIA_AT_CAMPAIGN_AND_ADGROUP",
	82: "HOTEL_LENGTH_OF_STAY_OVERLAPS_WITH_EXISTING_CRITERION",
	83: "HOTEL_ADVANCE_BOOKING_WINDOW_OVERLAPS_WITH_EXISTING_CRITERION",
	84: "FIELD_INCOMPATIBLE_WITH_NEGATIVE_TARGETING",
	85: "INVALID_WEBPAGE_CONDITION",
	86: "INVALID_WEBPAGE_CONDITION_URL",
	87: "WEBPAGE_CONDITION_URL_CANNOT_BE_EMPTY",
	88: "WEBPAGE_CONDITION_URL_UNSUPPORTED_PROTOCOL",
	89: "WEBPAGE_CONDITION_URL_CANNOT_BE_IP_ADDRESS",
	90: "WEBPAGE_CONDITION_URL_DOMAIN_NOT_CONSISTENT_WITH_CAMPAIGN_SETTING",
	91: "WEBPAGE_CONDITION_URL_CANNOT_BE_PUBLIC_SUFFIX",
	92: "WEBPAGE_CONDITION_URL_INVALID_PUBLIC_SUFFIX",
	93: "WEBPAGE_CONDITION_URL_VALUE_TRACK_VALUE_NOT_SUPPORTED",
	94: "WEBPAGE_CRITERION_URL_EQUALS_CAN_HAVE_ONLY_ONE_CONDITION",
	95: "WEBPAGE_CRITERION_NOT_SUPPORTED_ON_NON_DSA_AD_GROUP",
}
var CriterionErrorEnum_CriterionError_value = map[string]int32{
	"UNSPECIFIED":                                                       0,
	"UNKNOWN":                                                           1,
	"CONCRETE_TYPE_REQUIRED":                                            2,
	"INVALID_EXCLUDED_CATEGORY":                                         3,
	"INVALID_KEYWORD_TEXT":                                              4,
	"KEYWORD_TEXT_TOO_LONG":                                             5,
	"KEYWORD_HAS_TOO_MANY_WORDS":                                        6,
	"KEYWORD_HAS_INVALID_CHARS":                                         7,
	"INVALID_PLACEMENT_URL":                                             8,
	"INVALID_USER_LIST":                                                 9,
	"INVALID_USER_INTEREST":                                             10,
	"INVALID_FORMAT_FOR_PLACEMENT_URL":                                  11,
	"PLACEMENT_URL_IS_TOO_LONG":                                         12,
	"PLACEMENT_URL_HAS_ILLEGAL_CHAR":                                    13,
	"PLACEMENT_URL_HAS_MULTIPLE_SITES_IN_LINE":                          14,
	"PLACEMENT_IS_NOT_AVAILABLE_FOR_TARGETING_OR_EXCLUSION":             15,
	"INVALID_TOPIC_PATH":                                                16,
	"INVALID_YOUTUBE_CHANNEL_ID":                                        17,
	"INVALID_YOUTUBE_VIDEO_ID":                                          18,
	"YOUTUBE_VERTICAL_CHANNEL_DEPRECATED":                               19,
	"YOUTUBE_DEMOGRAPHIC_CHANNEL_DEPRECATED":                            20,
	"YOUTUBE_URL_UNSUPPORTED":                                           21,
	"CANNOT_EXCLUDE_CRITERIA_TYPE":                                      22,
	"CANNOT_ADD_CRITERIA_TYPE":                                          23,
	"INVALID_PRODUCT_FILTER":                                            24,
	"PRODUCT_FILTER_TOO_LONG":                                           25,
	"CANNOT_EXCLUDE_SIMILAR_USER_LIST":                                  26,
	"CANNOT_ADD_CLOSED_USER_LIST":                                       27,
	"CANNOT_ADD_DISPLAY_ONLY_LISTS_TO_SEARCH_ONLY_CAMPAIGNS":            28,
	"CANNOT_ADD_DISPLAY_ONLY_LISTS_TO_SEARCH_CAMPAIGNS":                 29,
	"CANNOT_ADD_DISPLAY_ONLY_LISTS_TO_SHOPPING_CAMPAIGNS":               30,
	"CANNOT_ADD_USER_INTERESTS_TO_SEARCH_CAMPAIGNS":                     31,
	"CANNOT_SET_BIDS_ON_CRITERION_TYPE_IN_SEARCH_CAMPAIGNS":             32,
	"CANNOT_ADD_URLS_TO_CRITERION_TYPE_FOR_CAMPAIGN_TYPE":               33,
	"INVALID_IP_ADDRESS":                                                34,
	"INVALID_IP_FORMAT":                                                 35,
	"INVALID_MOBILE_APP":                                                36,
	"INVALID_MOBILE_APP_CATEGORY":                                       37,
	"INVALID_CRITERION_ID":                                              38,
	"CANNOT_TARGET_CRITERION":                                           39,
	"CANNOT_TARGET_OBSOLETE_CRITERION":                                  40,
	"CRITERION_ID_AND_TYPE_MISMATCH":                                    41,
	"INVALID_PROXIMITY_RADIUS":                                          42,
	"INVALID_PROXIMITY_RADIUS_UNITS":                                    43,
	"INVALID_STREETADDRESS_LENGTH":                                      44,
	"INVALID_CITYNAME_LENGTH":                                           45,
	"INVALID_REGIONCODE_LENGTH":                                         46,
	"INVALID_REGIONNAME_LENGTH":                                         47,
	"INVALID_POSTALCODE_LENGTH":                                         48,
	"INVALID_COUNTRY_CODE":                                              49,
	"INVALID_LATITUDE":                                                  50,
	"INVALID_LONGITUDE":                                                 51,
	"PROXIMITY_GEOPOINT_AND_ADDRESS_BOTH_CANNOT_BE_NULL":                52,
	"INVALID_PROXIMITY_ADDRESS":                                         53,
	"INVALID_USER_DOMAIN_NAME":                                          54,
	"CRITERION_PARAMETER_TOO_LONG":                                      55,
	"AD_SCHEDULE_TIME_INTERVALS_OVERLAP":                                56,
	"AD_SCHEDULE_INTERVAL_CANNOT_SPAN_MULTIPLE_DAYS":                    57,
	"AD_SCHEDULE_INVALID_TIME_INTERVAL":                                 58,
	"AD_SCHEDULE_EXCEEDED_INTERVALS_PER_DAY_LIMIT":                      59,
	"AD_SCHEDULE_CRITERION_ID_MISMATCHING_FIELDS":                       60,
	"CANNOT_BID_MODIFY_CRITERION_TYPE":                                  61,
	"CANNOT_BID_MODIFY_CRITERION_CAMPAIGN_OPTED_OUT":                    62,
	"CANNOT_BID_MODIFY_NEGATIVE_CRITERION":                              63,
	"BID_MODIFIER_ALREADY_EXISTS":                                       64,
	"FEED_ID_NOT_ALLOWED":                                               65,
	"ACCOUNT_INELIGIBLE_FOR_CRITERIA_TYPE":                              66,
	"CRITERIA_TYPE_INVALID_FOR_BIDDING_STRATEGY":                        67,
	"CANNOT_EXCLUDE_CRITERION":                                          68,
	"CANNOT_REMOVE_CRITERION":                                           69,
	"PRODUCT_SCOPE_TOO_LONG":                                            70,
	"PRODUCT_SCOPE_TOO_MANY_DIMENSIONS":                                 71,
	"PRODUCT_PARTITION_TOO_LONG":                                        72,
	"PRODUCT_PARTITION_TOO_MANY_DIMENSIONS":                             73,
	"INVALID_PRODUCT_DIMENSION":                                         74,
	"INVALID_PRODUCT_DIMENSION_TYPE":                                    75,
	"INVALID_PRODUCT_BIDDING_CATEGORY":                                  76,
	"MISSING_SHOPPING_SETTING":                                          77,
	"INVALID_MATCHING_FUNCTION":                                         78,
	"LOCATION_FILTER_NOT_ALLOWED":                                       79,
	"LOCATION_FILTER_INVALID":                                           80,
	"CANNOT_ATTACH_CRITERIA_AT_CAMPAIGN_AND_ADGROUP":                    81,
	"HOTEL_LENGTH_OF_STAY_OVERLAPS_WITH_EXISTING_CRITERION":             82,
	"HOTEL_ADVANCE_BOOKING_WINDOW_OVERLAPS_WITH_EXISTING_CRITERION":     83,
	"FIELD_INCOMPATIBLE_WITH_NEGATIVE_TARGETING":                        84,
	"INVALID_WEBPAGE_CONDITION":                                         85,
	"INVALID_WEBPAGE_CONDITION_URL":                                     86,
	"WEBPAGE_CONDITION_URL_CANNOT_BE_EMPTY":                             87,
	"WEBPAGE_CONDITION_URL_UNSUPPORTED_PROTOCOL":                        88,
	"WEBPAGE_CONDITION_URL_CANNOT_BE_IP_ADDRESS":                        89,
	"WEBPAGE_CONDITION_URL_DOMAIN_NOT_CONSISTENT_WITH_CAMPAIGN_SETTING": 90,
	"WEBPAGE_CONDITION_URL_CANNOT_BE_PUBLIC_SUFFIX":                     91,
	"WEBPAGE_CONDITION_URL_INVALID_PUBLIC_SUFFIX":                       92,
	"WEBPAGE_CONDITION_URL_VALUE_TRACK_VALUE_NOT_SUPPORTED":             93,
	"WEBPAGE_CRITERION_URL_EQUALS_CAN_HAVE_ONLY_ONE_CONDITION":          94,
	"WEBPAGE_CRITERION_NOT_SUPPORTED_ON_NON_DSA_AD_GROUP":               95,
}

func (x CriterionErrorEnum_CriterionError) String() string {
	return proto.EnumName(CriterionErrorEnum_CriterionError_name, int32(x))
}
func (CriterionErrorEnum_CriterionError) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_criterion_error_0e10aab531921ef5, []int{0, 0}
}

// Container for enum describing possible criterion errors.
type CriterionErrorEnum struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CriterionErrorEnum) Reset()         { *m = CriterionErrorEnum{} }
func (m *CriterionErrorEnum) String() string { return proto.CompactTextString(m) }
func (*CriterionErrorEnum) ProtoMessage()    {}
func (*CriterionErrorEnum) Descriptor() ([]byte, []int) {
	return fileDescriptor_criterion_error_0e10aab531921ef5, []int{0}
}
func (m *CriterionErrorEnum) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CriterionErrorEnum.Unmarshal(m, b)
}
func (m *CriterionErrorEnum) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CriterionErrorEnum.Marshal(b, m, deterministic)
}
func (dst *CriterionErrorEnum) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CriterionErrorEnum.Merge(dst, src)
}
func (m *CriterionErrorEnum) XXX_Size() int {
	return xxx_messageInfo_CriterionErrorEnum.Size(m)
}
func (m *CriterionErrorEnum) XXX_DiscardUnknown() {
	xxx_messageInfo_CriterionErrorEnum.DiscardUnknown(m)
}

var xxx_messageInfo_CriterionErrorEnum proto.InternalMessageInfo

func init() {
	proto.RegisterType((*CriterionErrorEnum)(nil), "google.ads.googleads.v0.errors.CriterionErrorEnum")
	proto.RegisterEnum("google.ads.googleads.v0.errors.CriterionErrorEnum_CriterionError", CriterionErrorEnum_CriterionError_name, CriterionErrorEnum_CriterionError_value)
}

func init() {
	proto.RegisterFile("google/ads/googleads/v0/errors/criterion_error.proto", fileDescriptor_criterion_error_0e10aab531921ef5)
}

var fileDescriptor_criterion_error_0e10aab531921ef5 = []byte{
	// 1640 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x57, 0xef, 0x92, 0x1c, 0xb5,
	0x11, 0x8f, 0x4d, 0x02, 0x89, 0x9c, 0x80, 0x90, 0xff, 0xdb, 0xe7, 0xc3, 0x3e, 0x6c, 0x63, 0x8c,
	0xd9, 0x3b, 0xdb, 0xd8, 0xc0, 0x02, 0x49, 0xb4, 0xa3, 0xde, 0x59, 0xe5, 0x34, 0x92, 0x2c, 0x69,
	0xf6, 0x6e, 0xc9, 0x25, 0x2a, 0x87, 0x73, 0x5d, 0xb9, 0x0a, 0xbc, 0xd4, 0x1d, 0xe1, 0x81, 0xf2,
	0x31, 0x5f, 0xf3, 0x16, 0x3c, 0x4a, 0x3e, 0xe4, 0x19, 0x52, 0x3d, 0xb3, 0x9a, 0xd5, 0xec, 0xdd,
	0x61, 0x3e, 0xed, 0x56, 0xf7, 0xef, 0x27, 0x75, 0x4b, 0xbf, 0xee, 0xd6, 0x90, 0x4f, 0x0e, 0xe6,
	0xf3, 0x83, 0x6f, 0x5f, 0x6c, 0x3e, 0xdf, 0x3f, 0xda, 0x6c, 0xff, 0xe2, 0xbf, 0x1f, 0xb7, 0x36,
	0x5f, 0x1c, 0x1e, 0xce, 0x0f, 0x8f, 0x36, 0xbf, 0x39, 0x7c, 0xf9, 0xc3, 0x8b, 0xc3, 0x97, 0xf3,
	0x57, 0xb1, 0x31, 0x0c, 0xbe, 0x3f, 0x9c, 0xff, 0x30, 0x67, 0xeb, 0x2d, 0x74, 0xf0, 0x7c, 0xff,
	0x68, 0xd0, 0xb1, 0x06, 0x3f, 0x6e, 0x0d, 0x5a, 0xd6, 0xc6, 0x4f, 0x6b, 0x84, 0x15, 0x89, 0x09,
	0x68, 0x83, 0x57, 0xff, 0xfc, 0x6e, 0xe3, 0x3f, 0x6b, 0xe4, 0xed, 0xbe, 0x99, 0xbd, 0x43, 0xce,
	0xd5, 0xda, 0x5b, 0x28, 0xe4, 0x58, 0x82, 0xa0, 0xbf, 0x62, 0xe7, 0xc8, 0x5b, 0xb5, 0xde, 0xd6,
	0x66, 0x47, 0xd3, 0x33, 0xec, 0x1a, 0xb9, 0x54, 0x18, 0x5d, 0x38, 0x08, 0x10, 0xc3, 0xcc, 0x42,
	0x74, 0xf0, 0xac, 0x96, 0x0e, 0x04, 0x3d, 0xcb, 0x6e, 0x90, 0xab, 0x52, 0x4f, 0xb9, 0x92, 0x22,
	0xc2, 0x6e, 0xa1, 0x6a, 0x01, 0x22, 0x16, 0x3c, 0x40, 0x69, 0xdc, 0x8c, 0xbe, 0xc1, 0xae, 0x90,
	0x0b, 0xc9, 0xbd, 0x0d, 0xb3, 0x1d, 0xe3, 0x44, 0x0c, 0xb0, 0x1b, 0xe8, 0xaf, 0xd9, 0x55, 0x72,
	0x31, 0xb7, 0xc4, 0x60, 0x4c, 0x54, 0x46, 0x97, 0xf4, 0x37, 0x6c, 0x9d, 0x5c, 0x4b, 0xae, 0x09,
	0xf7, 0x8d, 0xa7, 0xe2, 0x7a, 0x16, 0xd1, 0xe2, 0xe9, 0x9b, 0xb8, 0x67, 0xee, 0x4f, 0x1b, 0x14,
	0x13, 0xee, 0x3c, 0x7d, 0x0b, 0x57, 0x4e, 0x26, 0xab, 0x78, 0x01, 0x15, 0xe8, 0x10, 0x6b, 0xa7,
	0xe8, 0x6f, 0xd9, 0x45, 0xf2, 0x6e, 0x72, 0xd5, 0x1e, 0x5c, 0x54, 0xd2, 0x07, 0xfa, 0xbb, 0x9c,
	0xd1, 0x98, 0xa5, 0x0e, 0xe0, 0xc0, 0x07, 0x4a, 0xd8, 0x6d, 0x72, 0x33, 0xb9, 0xc6, 0xc6, 0x55,
	0x3c, 0xe0, 0xcf, 0xca, 0xba, 0xe7, 0x30, 0xa2, 0x9e, 0x29, 0x4a, 0xbf, 0x4c, 0xe8, 0xf7, 0x6c,
	0x83, 0xac, 0xf7, 0xdd, 0x4d, 0xd8, 0x4a, 0x41, 0xc9, 0x55, 0x13, 0x36, 0xfd, 0x03, 0x7b, 0x40,
	0xee, 0x1d, 0xc7, 0x54, 0xb5, 0x0a, 0xd2, 0x2a, 0x88, 0x5e, 0x06, 0xc0, 0x4c, 0xa3, 0x92, 0x1a,
	0xe8, 0xdb, 0xec, 0x73, 0xf2, 0x64, 0x89, 0x96, 0x3e, 0x6a, 0x13, 0x22, 0x9f, 0x72, 0xa9, 0xf8,
	0x48, 0x41, 0x13, 0x62, 0xe0, 0xae, 0x84, 0x20, 0x75, 0x19, 0x8d, 0x6b, 0xaf, 0xc6, 0x4b, 0xa3,
	0xe9, 0x3b, 0xec, 0x12, 0x61, 0x29, 0xa3, 0x60, 0xac, 0x2c, 0xa2, 0xe5, 0x61, 0x42, 0x29, 0x9e,
	0x7a, 0xb2, 0xcf, 0x4c, 0x1d, 0xea, 0x11, 0x60, 0x68, 0x5a, 0x83, 0x8a, 0x52, 0xd0, 0x77, 0xd9,
	0x1a, 0xb9, 0xb2, 0xea, 0x9f, 0x4a, 0x01, 0x06, 0xbd, 0x8c, 0x7d, 0x40, 0xde, 0xef, 0xac, 0xe0,
	0x82, 0x2c, 0xda, 0xcc, 0x1a, 0xba, 0x00, 0xeb, 0x00, 0x45, 0x21, 0xe8, 0x79, 0x76, 0x9f, 0xdc,
	0x4d, 0x40, 0x01, 0x95, 0x29, 0x1d, 0xb7, 0x13, 0x59, 0x9c, 0x84, 0xbd, 0xc0, 0xae, 0x93, 0xcb,
	0x09, 0x8b, 0x27, 0x52, 0x6b, 0x5f, 0x5b, 0x6b, 0x1c, 0x3a, 0x2f, 0xb2, 0x9b, 0x64, 0xad, 0xe0,
	0x1a, 0x13, 0x5f, 0x08, 0x2f, 0x16, 0x4e, 0x06, 0x70, 0x92, 0x37, 0x22, 0xa5, 0x97, 0x30, 0xe2,
	0x05, 0x82, 0x0b, 0xb1, 0xe2, 0xbd, 0x8c, 0xaa, 0xee, 0x64, 0xe2, 0x8c, 0xa8, 0x8b, 0x10, 0xc7,
	0x52, 0x05, 0x70, 0xf4, 0x0a, 0x6e, 0xdc, 0xb7, 0x2d, 0x6f, 0xf3, 0x2a, 0x4a, 0x62, 0x65, 0x63,
	0x2f, 0x2b, 0xa9, 0xb8, 0xcb, 0x34, 0x75, 0x8d, 0xbd, 0x47, 0xae, 0xe7, 0x9b, 0x2b, 0xe3, 0x21,
	0x17, 0xdd, 0x75, 0x36, 0x24, 0x4f, 0x33, 0x80, 0x90, 0xde, 0x2a, 0x3e, 0x8b, 0x46, 0xab, 0x59,
	0x83, 0x40, 0x01, 0x45, 0x0f, 0xdc, 0x15, 0x93, 0xd6, 0x58, 0xf0, 0xca, 0x72, 0x59, 0x6a, 0x4f,
	0xd7, 0xd8, 0x13, 0xf2, 0xf0, 0x97, 0x72, 0x97, 0xb4, 0x1b, 0xec, 0x53, 0xf2, 0xf8, 0xf5, 0xb4,
	0x89, 0xb1, 0x16, 0x65, 0xb3, 0x24, 0xae, 0xb3, 0x87, 0xe4, 0xe3, 0x8c, 0xd8, 0xab, 0x91, 0x93,
	0xf7, 0x7a, 0x0f, 0x15, 0xba, 0xa0, 0x78, 0x08, 0x71, 0x24, 0x85, 0x8f, 0x46, 0xa7, 0x4b, 0x30,
	0xba, 0x6d, 0x24, 0x52, 0x1f, 0xa7, 0xde, 0x5c, 0x09, 0xb3, 0x76, 0xaa, 0xd9, 0x63, 0x85, 0x8a,
	0x1a, 0x4f, 0xa4, 0xf6, 0x4a, 0x6f, 0xe5, 0xd2, 0x96, 0x16, 0xc9, 0x0e, 0xbc, 0xa7, 0x1b, 0x79,
	0xd9, 0x4b, 0xbb, 0xa8, 0x63, 0xfa, 0x7e, 0x0e, 0xaf, 0xcc, 0x48, 0x2a, 0x88, 0xdc, 0x5a, 0x7a,
	0x1b, 0xaf, 0xee, 0xb8, 0x7d, 0xd9, 0xd5, 0xee, 0xe4, 0x5d, 0x6d, 0x19, 0x95, 0x14, 0xf4, 0x2e,
	0x0a, 0x67, 0x11, 0x7a, 0x5b, 0x7f, 0x4b, 0x3f, 0xfd, 0x20, 0x13, 0xce, 0xc2, 0x69, 0x46, 0xde,
	0x28, 0xec, 0xaa, 0x4b, 0xd4, 0x3d, 0x6c, 0x16, 0xf9, 0xa2, 0x91, 0x6b, 0xd1, 0x66, 0x5c, 0x49,
	0x5f, 0xf1, 0x50, 0x4c, 0xe8, 0x87, 0x79, 0x2d, 0x5a, 0x67, 0x76, 0x65, 0x25, 0xc3, 0x2c, 0x3a,
	0x2e, 0x64, 0xed, 0xe9, 0x7d, 0x5c, 0xe1, 0x34, 0x6f, 0xac, 0xb5, 0x0c, 0x9e, 0x7e, 0x84, 0xd5,
	0x93, 0x30, 0x3e, 0x38, 0x80, 0xb0, 0x38, 0xad, 0xa8, 0x40, 0x97, 0x61, 0x42, 0x1f, 0x60, 0x2a,
	0x5d, 0x92, 0x32, 0xcc, 0x34, 0xaf, 0x20, 0x39, 0x3f, 0xce, 0xdb, 0xbe, 0x83, 0x52, 0x1a, 0x5d,
	0x18, 0xd1, 0xb9, 0x07, 0xc7, 0xdd, 0x39, 0x7b, 0x33, 0x77, 0x5b, 0xe3, 0x03, 0x57, 0x39, 0x7b,
	0xab, 0x77, 0xbc, 0xa6, 0xd6, 0xc1, 0xcd, 0x22, 0x02, 0xe8, 0x43, 0x76, 0x81, 0xd0, 0xe4, 0x51,
	0x3c, 0xc8, 0x50, 0x0b, 0xa0, 0x8f, 0xf2, 0xeb, 0xc5, 0x12, 0x6d, 0xcd, 0x8f, 0xd9, 0x53, 0xf2,
	0x68, 0x99, 0x7e, 0x09, 0xc6, 0x1a, 0xa9, 0x43, 0x73, 0x9c, 0x29, 0xd7, 0x91, 0x09, 0xa8, 0xbc,
	0xe6, 0x4a, 0x46, 0x10, 0x75, 0xad, 0x14, 0xfd, 0xa4, 0x17, 0x5d, 0xc7, 0x4f, 0x62, 0x7a, 0x92,
	0x9f, 0x7d, 0x53, 0x08, 0xc2, 0x54, 0x5c, 0xea, 0x88, 0x19, 0xd2, 0xa7, 0x4d, 0x57, 0xea, 0x6e,
	0xcf, 0x72, 0xc7, 0x2b, 0xe8, 0xb5, 0x8f, 0x4f, 0xd9, 0x5d, 0xb2, 0xc1, 0x45, 0xf4, 0xc5, 0x04,
	0x44, 0xad, 0x20, 0x06, 0x59, 0x41, 0x5b, 0x4c, 0x53, 0xae, 0x7c, 0x34, 0x53, 0x70, 0x8a, 0x5b,
	0xfa, 0x19, 0x7b, 0x44, 0x06, 0x39, 0x2e, 0x41, 0x52, 0xbc, 0xde, 0x72, 0xbd, 0x9c, 0x0f, 0x82,
	0xcf, 0x3c, 0xfd, 0x9c, 0xdd, 0x21, 0xb7, 0xfa, 0x9c, 0x45, 0x9f, 0xcf, 0xf7, 0xa0, 0x43, 0xb6,
	0x45, 0x1e, 0xe4, 0x30, 0xd8, 0x2d, 0x00, 0x70, 0x70, 0x2f, 0xc3, 0xb0, 0x98, 0x18, 0xc7, 0xb6,
	0x50, 0xc9, 0x40, 0xbf, 0x60, 0x9b, 0xe4, 0xa3, 0x9c, 0xd1, 0x13, 0x68, 0xd2, 0x25, 0x36, 0x8d,
	0xb1, 0x04, 0x25, 0x3c, 0xfd, 0x32, 0xd3, 0xfa, 0xa8, 0x29, 0x23, 0x21, 0xc7, 0xb3, 0x95, 0x12,
	0xa6, 0x5f, 0x61, 0x8e, 0x3f, 0x87, 0xea, 0xea, 0xdb, 0xd8, 0x00, 0x22, 0x9a, 0x3a, 0xd0, 0x3f,
	0xb2, 0x7b, 0xe4, 0xf6, 0x71, 0x8e, 0x86, 0x92, 0x07, 0x39, 0xcd, 0x2b, 0xe9, 0x4f, 0x58, 0xc7,
	0x1d, 0x44, 0x82, 0x8b, 0x5c, 0x39, 0xe0, 0x62, 0x16, 0x61, 0x17, 0x9b, 0x1d, 0xfd, 0x33, 0xbb,
	0x4c, 0xce, 0x8f, 0x01, 0xf3, 0x16, 0xed, 0x00, 0x55, 0xca, 0xec, 0x80, 0xa0, 0x1c, 0xf7, 0xe0,
	0x45, 0xa3, 0xbd, 0x28, 0x35, 0x28, 0x59, 0xca, 0x34, 0x56, 0xfb, 0x53, 0x64, 0xc4, 0x06, 0xe4,
	0x7e, 0xcf, 0x14, 0xb3, 0xd7, 0x02, 0x06, 0x28, 0xf0, 0x54, 0x7c, 0x70, 0xd8, 0x3c, 0x66, 0xb4,
	0xc8, 0x66, 0xd2, 0xca, 0xd4, 0x32, 0x9a, 0x8a, 0xac, 0x7d, 0x38, 0xa8, 0x4c, 0x2f, 0x1d, 0xc0,
	0x81, 0x95, 0x86, 0x92, 0x2f, 0x8c, 0x85, 0xa5, 0xa8, 0xc6, 0x78, 0xf1, 0xc7, 0x7d, 0xcd, 0xa3,
	0x49, 0xc8, 0x0a, 0x34, 0x8e, 0x7e, 0x4f, 0x4b, 0x9c, 0xf1, 0x09, 0x66, 0xb9, 0x0b, 0x32, 0x34,
	0x77, 0x91, 0x96, 0x99, 0xb0, 0x0f, 0xc9, 0x9d, 0x93, 0xfd, 0xab, 0x4b, 0xc9, 0x95, 0x2a, 0x69,
	0x28, 0x9d, 0x9f, 0xfe, 0x65, 0xa5, 0x07, 0xf5, 0xdd, 0xed, 0xd9, 0x6d, 0xe7, 0x6f, 0xab, 0x84,
	0x49, 0x27, 0xd6, 0x35, 0x5b, 0x85, 0x27, 0x56, 0x49, 0xef, 0x9b, 0x73, 0x4c, 0xb3, 0xc9, 0x43,
	0xc0, 0xa7, 0x0d, 0xad, 0xf2, 0x30, 0x96, 0x22, 0xac, 0x75, 0x81, 0x91, 0x53, 0x8d, 0x12, 0x50,
	0xa6, 0xe0, 0x4d, 0x1e, 0x8b, 0x49, 0x9e, 0xdf, 0xb4, 0xc1, 0x13, 0x5f, 0x05, 0x2c, 0xd6, 0xa3,
	0x36, 0x93, 0x27, 0x0f, 0x81, 0xe3, 0x94, 0x4a, 0x57, 0xcd, 0xc3, 0x52, 0x9c, 0x6d, 0x47, 0x29,
	0x9d, 0xa9, 0x2d, 0x7d, 0x86, 0x73, 0x6f, 0x62, 0x02, 0xa8, 0x45, 0x3b, 0x8b, 0x66, 0x1c, 0x7d,
	0xc0, 0x29, 0xdb, 0x16, 0xb6, 0x8f, 0x3b, 0x32, 0x4c, 0x5a, 0x09, 0x36, 0x89, 0x76, 0x17, 0xec,
	0x18, 0x27, 0x5f, 0xb5, 0x54, 0x2e, 0xa6, 0x5c, 0x17, 0x10, 0x47, 0xc6, 0x6c, 0x23, 0x68, 0x47,
	0x6a, 0x61, 0x76, 0x5e, 0xbf, 0x84, 0x47, 0x39, 0x36, 0x25, 0x18, 0xa5, 0x2e, 0x4c, 0x65, 0x79,
	0x68, 0x84, 0xdb, 0xe0, 0xbb, 0x12, 0xe9, 0x5e, 0x86, 0x34, 0xe4, 0xc7, 0xb7, 0x03, 0x23, 0xcb,
	0x4b, 0x88, 0x85, 0xd1, 0xa2, 0xb9, 0x78, 0x5a, 0xb3, 0x5b, 0xe4, 0xc6, 0xa9, 0xee, 0xe6, 0xe9,
	0x3b, 0x45, 0xc9, 0x9c, 0xe8, 0xca, 0xfa, 0x2a, 0x54, 0x36, 0xcc, 0xe8, 0x0e, 0x06, 0x77, 0x32,
	0x34, 0x7b, 0xd8, 0xa1, 0x12, 0x82, 0x29, 0x8c, 0xa2, 0xbb, 0xa7, 0xe3, 0x97, 0x4b, 0x67, 0x63,
	0x7e, 0xc6, 0x80, 0xf0, 0x93, 0xf1, 0xa9, 0x45, 0x9b, 0x80, 0x0e, 0x2f, 0x7d, 0xc0, 0x67, 0x73,
	0x73, 0x2a, 0xdd, 0x2d, 0x26, 0x49, 0x7d, 0x8d, 0x8f, 0x9d, 0xd7, 0x6d, 0x6b, 0xeb, 0x91, 0x92,
	0x45, 0xf4, 0xf5, 0x78, 0x2c, 0x77, 0xe9, 0x5f, 0xb1, 0x3d, 0x9e, 0x4c, 0xe9, 0xf4, 0xdd, 0x23,
	0xec, 0xa1, 0x4a, 0x4e, 0x26, 0x4c, 0xb9, 0xaa, 0x21, 0x06, 0xc7, 0x8b, 0xed, 0xc5, 0xff, 0xa6,
	0xd9, 0x77, 0xef, 0xde, 0xbf, 0xb1, 0x2f, 0xc9, 0x67, 0x1d, 0xb5, 0xeb, 0x94, 0x48, 0x85, 0x67,
	0x35, 0xb6, 0xee, 0x82, 0xeb, 0x38, 0xe1, 0x53, 0x68, 0xdf, 0x75, 0x46, 0xe7, 0x37, 0xfa, 0x77,
	0x7c, 0x5b, 0x1d, 0x67, 0xf7, 0xb6, 0x88, 0x8d, 0x41, 0x47, 0xe1, 0x79, 0xe4, 0x22, 0xb6, 0xba,
	0x8e, 0xa3, 0xff, 0x9d, 0x21, 0x1b, 0xdf, 0xcc, 0xbf, 0x1b, 0xfc, 0xfc, 0x37, 0xe7, 0xe8, 0x7c,
	0xff, 0xcb, 0xd2, 0xe2, 0x87, 0xaa, 0x3d, 0xf3, 0xb5, 0x58, 0xd0, 0x0e, 0xe6, 0xdf, 0x3e, 0x7f,
	0x75, 0x30, 0x98, 0x1f, 0x1e, 0x6c, 0x1e, 0xbc, 0x78, 0xd5, 0x7c, 0xc6, 0xa6, 0x0f, 0xde, 0xef,
	0x5f, 0x1e, 0x9d, 0xf6, 0xfd, 0xfb, 0x45, 0xfb, 0xf3, 0xaf, 0xb3, 0x6f, 0x94, 0x9c, 0xff, 0xfb,
	0xec, 0x7a, 0xd9, 0x2e, 0xc6, 0xf7, 0x8f, 0x06, 0xed, 0x5f, 0xfc, 0x37, 0xdd, 0x1a, 0x34, 0x5b,
	0x1e, 0xfd, 0x94, 0x00, 0x7b, 0x7c, 0xff, 0x68, 0xaf, 0x03, 0xec, 0x4d, 0xb7, 0xf6, 0x5a, 0xc0,
	0x7f, 0xcf, 0x6e, 0xb4, 0xd6, 0xe1, 0x90, 0xef, 0x1f, 0x0d, 0x87, 0x1d, 0x64, 0x38, 0x9c, 0x6e,
	0x0d, 0x87, 0x2d, 0xe8, 0x1f, 0x6f, 0x36, 0xd1, 0x3d, 0xfe, 0x7f, 0x00, 0x00, 0x00, 0xff, 0xff,
	0x80, 0xa6, 0xb7, 0xd7, 0x9c, 0x0f, 0x00, 0x00,
}
