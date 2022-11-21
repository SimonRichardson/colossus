package selectors

// UserType represents the type of user (admin, phone, email, etc)
type UserType uint64

func (t UserType) XOR(x UserType) UserType {
	return t ^ x
}

func (t UserType) Val() uint64 {
	return uint64(t)
}

const (
	UserTypeInvalid              UserType = 1 << 0
	UserTypePhone                UserType = 1 << 1
	UserTypePromoter             UserType = 1 << 2
	UserTypeAgent                UserType = 1 << 3
	UserTypeManager              UserType = 1 << 4
	UserTypeArtist               UserType = 1 << 5
	UserTypeAdmin                UserType = 1 << 6
	UserTypeRoot                 UserType = 1 << 7
	UserTypeRep                  UserType = 1 << 8
	UserTypePartialPhone         UserType = 1 << 9
	UserTypeSocialFacebookLegacy UserType = 1 << 11
	UserTypeSocialGoogle         UserType = 1 << 12
	UserTypeBlockedForTransfer   UserType = 1 << 13
	UserTypeVIP                  UserType = 1 << 14
	UserTypeUser                 UserType = 1 << 15
)
