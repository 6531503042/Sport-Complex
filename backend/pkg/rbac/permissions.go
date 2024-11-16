package rbac

import "main/modules/auth"

// RolePermissions maps role codes to their permissions
var RolePermissions = map[int][]string{
    auth.RoleUser: {
        auth.PermissionReadUser,
    },
    auth.RoleAdmin: {
        auth.PermissionReadUser,
        auth.PermissionCreateUser,
        auth.PermissionUpdateUser,
        auth.PermissionDeleteUser,
        auth.PermissionAccessDashboard,
        auth.PermissionManageBookings,
    },
}

// HasPermission checks if a role has a specific permission
func HasPermission(roleCode int, permission string) bool {
    permissions, exists := RolePermissions[roleCode]
    if !exists {
        return false
    }

    // Admin has access to everything
    if roleCode == auth.RoleAdmin {
        return true
    }

    // For other roles, check specific permissions
    for _, p := range permissions {
        if p == permission {
            return true
        }
    }
    return false
} 