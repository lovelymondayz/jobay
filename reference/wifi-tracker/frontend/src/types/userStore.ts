export interface User {
  name: string;
  email: string;
  password: string;
  role_id: number;
  client_id: string;
  // user_id: string;
  created_by: string;
  updated_by: string;
}

export interface UserFormProps {
  user?: User | null;
  onSubmit: (data: User) => void;
  onCancel: () => void;
}
