"use client";

import { useState, useEffect } from "react";
import { Users, Clock, Shield, User, Mail, Phone } from "lucide-react";
import { getUsers } from "@/lib/api/console";
import type { UserSummary } from "@/lib/api/console";

export default function ConsolePeople() {
  const [users, setUsers] = useState<UserSummary[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true);
    getUsers(1, 50)
      .then((result) => setUsers(result.users))
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-text-charcoal">People</h1>
        <p className="text-text-muted mt-1">Team members, reporters, and platform users.</p>
      </div>

      {/* Summary */}
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
        <div className="bg-background rounded-xl border border-background-stone p-5 flex items-center gap-4">
          <div className="p-3 rounded-lg bg-forest/10 text-forest">
            <Users className="h-5 w-5" />
          </div>
          <div>
            <p className="text-2xl font-bold text-text-charcoal">{users.length}</p>
            <p className="text-xs text-text-muted">Total Users</p>
          </div>
        </div>
        <div className="bg-background rounded-xl border border-background-stone p-5 flex items-center gap-4">
          <div className="p-3 rounded-lg bg-status-important/10 text-status-important">
            <Shield className="h-5 w-5" />
          </div>
          <div>
            <p className="text-2xl font-bold text-text-charcoal">{users.filter((u) => u.is_platform_admin).length}</p>
            <p className="text-xs text-text-muted">Admins</p>
          </div>
        </div>
        <div className="bg-background rounded-xl border border-background-stone p-5 flex items-center gap-4">
          <div className="p-3 rounded-lg bg-sky/10 text-sky">
            <User className="h-5 w-5" />
          </div>
          <div>
            <p className="text-2xl font-bold text-text-charcoal">{users.filter((u) => !u.is_platform_admin).length}</p>
            <p className="text-xs text-text-muted">Regular Users</p>
          </div>
        </div>
      </div>

      {/* User list */}
      {loading ? (
        <div className="rounded-xl border border-background-stone bg-background p-8 text-center">
          <Clock className="h-6 w-6 text-text-muted mx-auto mb-2 animate-pulse" />
          <p className="text-text-muted text-sm">Loading users...</p>
        </div>
      ) : users.length === 0 ? (
        <div className="rounded-xl border border-background-stone bg-background p-8 text-center">
          <Users className="h-10 w-10 text-text-muted mx-auto mb-3" />
          <h2 className="text-lg font-semibold text-text-charcoal">No users found</h2>
          <p className="text-sm text-text-muted mt-2">Users will appear here once they register.</p>
        </div>
      ) : (
        <div className="bg-background rounded-xl border border-background-stone overflow-hidden">
          <table className="w-full text-left text-sm">
            <thead className="bg-background-mist border-b border-background-stone text-text-muted">
              <tr>
                <th className="p-4 font-medium">Name</th>
                <th className="p-4 font-medium">Email</th>
                <th className="p-4 font-medium">Phone</th>
                <th className="p-4 font-medium">Language</th>
                <th className="p-4 font-medium">Role</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-background-stone">
              {users.map((user) => (
                <tr key={user.id} className="hover:bg-background-mist/50">
                  <td className="p-4">
                    <div className="flex items-center gap-3">
                      <div className="h-8 w-8 rounded-full bg-forest text-white flex items-center justify-center text-xs font-medium">
                        {user.display_name?.charAt(0)?.toUpperCase() || "?"}
                      </div>
                      <span className="text-sm font-medium text-text-charcoal">{user.display_name}</span>
                    </div>
                  </td>
                  <td className="p-4 text-text-muted">
                    {user.email ? (
                      <span className="flex items-center gap-1"><Mail className="h-3 w-3" /> {user.email}</span>
                    ) : (
                      <span className="text-text-muted/50">—</span>
                    )}
                  </td>
                  <td className="p-4 text-text-muted">
                    {user.phone_number ? (
                      <span className="flex items-center gap-1"><Phone className="h-3 w-3" /> {user.phone_number}</span>
                    ) : (
                      <span className="text-text-muted/50">—</span>
                    )}
                  </td>
                  <td className="p-4 text-text-muted uppercase text-xs">{user.language}</td>
                  <td className="p-4">
                    {user.is_platform_admin ? (
                      <span className="inline-flex items-center gap-1 px-2 py-1 rounded-full text-xs font-medium bg-status-important/10 text-status-important border border-status-important/20">
                        <Shield className="h-3 w-3" /> Admin
                      </span>
                    ) : (
                      <span className="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-background-stone text-text-muted">
                        User
                      </span>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
