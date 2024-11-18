export interface IUserRole {
    role_title: string;
    role_code: number;
}

export interface IUser {
    _id: string;
    email: string;
    name: string;
    password?: string;
    created_at: string;
    updated_at: string;
    user_roles: IUserRole[];
}
