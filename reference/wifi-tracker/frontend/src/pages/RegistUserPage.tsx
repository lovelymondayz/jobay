import { useState } from "react";
import type { User, UserFormProps } from "@/types/userStore";
import { v4 as uuidv4 } from "uuid"; // install with: npm install uuid

// Import hooks
import {
  useUsers,
  useCreateUser,
  useUpdateUser,
  useDeleteUser,
} from "@/hooks/useActiveUser";

// User Form
import { useClients } from "@/hooks/useActiveClients";

const UserForm = ({ user = null, onSubmit, onCancel }: UserFormProps) => {
  const { data: clients, isLoading: loadingClients } = useClients();

  const [formData, setFormData] = useState<User>(
    user ?? {
      name: "",
      email: "",
      password: "",
      role_id: 0,
      client_id: "",
      // user_id: "",
      created_by: "",
      updated_by: "",
    }
  );

  const handleSubmit = () => {
    const payload = {
      ...formData,
      // user_id: user ? formData.user_id : uuidv4(),
      client_id: user ? formData.client_id : uuidv4(),
      created_by: user ? formData.created_by : "admin",
      updated_by: user ? formData.updated_by : "admin",
    };
    onSubmit(payload);
  };

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: name === "role_id" ? Number(value) : value,
    }));
  };

  return (
    <div className="space-y-4 p-4 border rounded-lg bg-gray-50">
      <div className="grid grid-cols-2 gap-4">
        {/* Client ID dropdown */}
        <div>
          <label className="block text-sm font-medium mb-1">Client</label>
          {loadingClients ? (
            <p className="text-gray-500 text-sm">Loading clients...</p>
          ) : (
            <select
              name="client_id"
              value={formData.client_id}
              onChange={handleChange}
              className="w-full p-2 border rounded-md"
            >
              <option value="">Select Client</option>
              {clients?.map((client: any) => (
                <option key={client.id} value={client.id}>
                  {client.name}
                </option>
              ))}
            </select>
          )}
        </div>

        {/* User ID (read-only) */}
        {/* <div>
          <label className="block text-sm font-medium mb-1">User ID</label>
          <input
            type="text"
            name="user_id"
            disabled
            value={formData.user_id}
            className="w-full p-2 border rounded-md"
          />
        </div> */}

        {/* Other inputs... */}
        <div>
          <label className="block text-sm font-medium mb-1">Name</label>
          <input
            type="text"
            name="name"
            value={formData.name}
            onChange={handleChange}
            className="w-full p-2 border rounded-md"
          />
        </div>
        <div>
          <label className="block text-sm font-medium mb-1">Email</label>
          <input
            type="text"
            name="email"
            value={formData.email}
            onChange={handleChange}
            placeholder="mack.jonson@gmail.com"
            className="w-full p-2 border rounded-md"
          />
        </div>
        <div>
          <label className="block text-sm font-medium mb-1">Password</label>
          <input
            type="password"
            name="password"
            value={formData.password}
            onChange={handleChange}
            className="w-full p-2 border rounded-md"
          />
        </div>
        <div>
          <label className="block text-sm font-medium mb-1">Role</label>
          <select
            name="role_id"
            value={formData.role_id}
            onChange={handleChange}
            className="w-full p-2 border rounded-md"
          >
            <option value={1}>Admin</option>
            <option value={2}>User</option>
            <option value={3}>Manager</option>
          </select>
        </div>
      </div>

      {/* Buttons */}
      <div className="flex gap-2">
        <button
          onClick={handleSubmit}
          className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
        >
          {user ? "Update" : "Create"}
        </button>
        <button
          onClick={onCancel}
          className="px-4 py-2 bg-gray-400 text-white rounded-md hover:bg-gray-500"
        >
          Cancel
        </button>
      </div>
    </div>
  );
};

// CRUD Page
const RegistUserPage = () => {
  const [showForm, setShowForm] = useState(false);
  const [editingUser, setEditingUser] = useState<User | null>(null);

  const { data: users, isLoading, error } = useUsers();
  const createMutation = useCreateUser();
  const updateMutation = useUpdateUser();
  const deleteMutation = useDeleteUser();

  const handleCreate = (data: User) => {
    createMutation.mutate(data, {
      onSuccess: () => {
        setShowForm(false);
      },
    });
  };

  const handleUpdate = (data: User) => {
    if (!editingUser) return;
    updateMutation.mutate(
      { id: editingUser.client_id, data },
      {
        onSuccess: () => {
          setEditingUser(null);
          setShowForm(false);
        },
      }
    );
  };

  const handleDelete = (user: User) => {
    if (window.confirm(`Delete ${user.name}?`)) {
      deleteMutation.mutate(user.client_id);
    }
  };

  const startEdit = (user: User) => {
    setEditingUser(user);
    setShowForm(true);
  };

  const cancelForm = () => {
    setShowForm(false);
    setEditingUser(null);
  };

  if (isLoading) return <div className="p-4">Loading users...</div>;
  if (error)
    return <div className="p-4 text-red-600">Error: {error.message}</div>;

  return (
    <div className="p-6 max-w-6xl mx-auto">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Users</h1>
        <button
          onClick={() => {
            setEditingUser(null);
            setShowForm(true);
          }}
          className="px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700"
        >
          Add User
        </button>
      </div>

      {showForm && (
        <div className="mb-6">
          <h2 className="text-lg font-semibold mb-3">
            {editingUser ? "Edit User" : "Add New User"}
          </h2>
          <UserForm
            user={editingUser}
            onSubmit={editingUser ? handleUpdate : handleCreate}
            onCancel={cancelForm}
          />
        </div>
      )}

      <div className="grid gap-4">
        {Array.isArray(users) &&
          users.map((user) => (
            <div
              key={user.client_id}
              className="border rounded-lg p-4 bg-white shadow-sm"
            >
              <div className="flex justify-between items-start">
                <div className="grid grid-cols-2 md:grid-cols-4 gap-4 flex-1">
                  <div>
                    <span className="text-sm text-gray-600">Name:</span>
                    <p className="font-medium">{user.name}</p>
                  </div>
                  <div>
                    <span className="text-sm text-gray-600">Client ID:</span>
                    <p className="text-xs font-mono">{user.client_id}</p>
                  </div>
                  <div>
                    <span className="text-sm text-gray-600">Email:</span>
                    <p className="font-medium">{user.email}</p>
                  </div>
                  <div>
                    <span className="text-sm text-gray-600">Role:</span>
                    <p className="font-medium">{user.role_id}</p>
                  </div>
                </div>

                <div className="flex gap-2">
                  <button
                    onClick={() => startEdit(user)}
                    className="px-3 py-1 bg-yellow-500 text-white rounded-md hover:bg-yellow-600"
                  >
                    Edit
                  </button>
                  <button
                    onClick={() => handleDelete(user)}
                    className="px-3 py-1 bg-red-600 text-white rounded-md hover:bg-red-700"
                  >
                    Delete
                  </button>
                </div>
              </div>
            </div>
          ))}
      </div>

      {users?.length === 0 && (
        <div className="text-center py-8 text-gray-500">
          No users found. Add your first user!
        </div>
      )}
    </div>
  );
};

export default RegistUserPage;
