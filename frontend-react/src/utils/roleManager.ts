import { jwtDecode } from 'jwt-decode';

export enum UserRole {
  OWNER = 'owner',
  FARM_MANAGER = 'farm_manager'
}

export interface UserPermissions {
  viewHorseData: boolean;
  editHorseData: boolean;
  exportData: boolean;
  manageFarm: boolean;
}

const ROLE_PERMISSIONS: Record<UserRole, UserPermissions> = {
  [UserRole.OWNER]: {
    viewHorseData: true,
    editHorseData: false,
    exportData: true,
    manageFarm: false
  },
  [UserRole.FARM_MANAGER]: {
    viewHorseData: true,
    editHorseData: true,
    exportData: true,
    manageFarm: true
  }
};

export class RoleManager {
  private static instance: RoleManager;
  private currentRole: UserRole = UserRole.OWNER;

  private constructor() {}

  public static getInstance(): RoleManager {
    if (!RoleManager.instance) {
      RoleManager.instance = new RoleManager();
    }
    return RoleManager.instance;
  }

  public setRoleFromToken(token: string): void {
    try {
      const decoded = jwtDecode<{ role?: UserRole }>(token);
      this.currentRole = decoded.role || UserRole.OWNER;
    } catch (error) {
      console.warn('Invalid token, defaulting to OWNER role');
      this.currentRole = UserRole.OWNER;
    }
  }

  public getCurrentRole(): UserRole {
    return this.currentRole;
  }

  public can(permission: keyof UserPermissions): boolean {
    return ROLE_PERMISSIONS[this.currentRole][permission];
  }

  public hasRole(role: UserRole): boolean {
    return this.currentRole === role || 
           this.getRoleHierarchy().includes(role);
  }

  private getRoleHierarchy(): UserRole[] {
    const roleOrder = [
      UserRole.OWNER, 
      UserRole.FARM_MANAGER
    ];
    
    const currentIndex = roleOrder.indexOf(this.currentRole);
    return roleOrder.slice(0, currentIndex + 1);
  }

  public getPermissions(): UserPermissions {
    return ROLE_PERMISSIONS[this.currentRole];
  }
}

// Utility hook for React components
export function useRoleManager() {
  return RoleManager.getInstance();
}

// Example of a permission-based component wrapper
export function withRoleCheck<P>(
  WrappedComponent: React.ComponentType<P>, 
  requiredPermission?: keyof UserPermissions
) {
  return (props: P) => {
    const roleManager = useRoleManager();
    
    if (requiredPermission && !roleManager.can(requiredPermission)) {
      return (
        <div>
          <h2>Access Denied</h2>
          <p>You do not have permission to view this page.</p>
        </div>
      );
    }

    return <WrappedComponent {...props} />;
  };
}
