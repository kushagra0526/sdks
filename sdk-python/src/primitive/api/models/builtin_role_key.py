from enum import StrEnum

class BuiltinRoleKey(StrEnum):
    ADMIN = "admin"
    MEMBER = "member"
    OPERATOR = "operator"
    OWNER = "owner"
    RESTRICTED = "restricted"

    def __str__(self) -> str:
        return str(self.value)
