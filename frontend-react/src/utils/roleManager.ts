export enum UserRole {
  GUEST = 'guest',
  USER = 'user',
  ADMIN = 'admin'
}

export interface UserPermission {
  userId: string;
  role: UserRole;
}

export class RoleManager {
  private static userRoles: Map<string, UserRole> = new Map();

  static assignRole(userId: string, role: UserRole): void {
      this.userRoles.set(userId, role);
  }

  static getUserRole(userId: string): UserRole | undefined {
      return this.userRoles.get(userId);
  }

  static hasPermission(
      userId: string, 
      requiredRole: UserRole
  ): boolean {
      const roleHierarchy: Record<UserRole, number> = {
          [UserRole.GUEST]: 1,
          [UserRole.USER]: 2,
          [UserRole.ADMIN]: 3
      };

      const userRole = this.getUserRole(userId);
      
      return userRole 
          ? roleHierarchy[userRole] >= roleHierarchy[requiredRole]
          : false;
  }

  static getAllUsersWithRole(role: UserRole): string[] {
      return Array.from(this.userRoles.entries())
          .filter(([_, userRole]) => userRole === role)
          .map(([userId, _]) => userId);
  }
}

export default RoleManager;