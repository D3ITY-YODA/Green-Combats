import { Plus } from "lucide-react";

const mockUpdates = [
  { status: "Published", type: "Flood warning", place: "Lower Valley", updated: "10:00" },
  { status: "Review", type: "Water outlook", place: "East Ward", updated: "09:30" },
  { status: "Draft", type: "Seasonal info", place: "North Basin", updated: "Yesterday" },
];

export default function ConsoleUpdates() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-text-charcoal">Updates</h1>
        <button className="flex items-center gap-2 px-4 py-2 rounded-lg bg-forest text-white text-sm font-medium hover:bg-forest-deep">
          <Plus className="h-4 w-4" /> Create update
        </button>
      </div>

      <div className="bg-background rounded-xl border border-background-stone overflow-hidden">
        <table className="w-full text-left text-sm">
          <thead className="bg-background-mist border-b border-background-stone text-text-muted">
            <tr>
              <th className="p-4 font-medium">Status</th>
              <th className="p-4 font-medium">Type</th>
              <th className="p-4 font-medium">Place</th>
              <th className="p-4 font-medium">Updated</th>
              <th className="p-4 font-medium text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-background-stone">
            {mockUpdates.map((update, i) => (
              <tr key={i} className="hover:bg-background-mist/50">
                <td className="p-4">
                  <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
                    update.status === "Published" ? "bg-status-normal/10 text-status-normal" :
                    update.status === "Review" ? "bg-status-watch/10 text-status-watch" :
                    "bg-background-stone text-text-muted"
                  }`}>
                    {update.status}
                  </span>
                </td>
                <td className="p-4 text-text-charcoal font-medium">{update.type}</td>
                <td className="p-4 text-text-muted">{update.place}</td>
                <td className="p-4 text-text-muted">{update.updated}</td>
                <td className="p-4 text-right">
                  <button className="text-forest hover:underline text-xs font-medium">Edit</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
