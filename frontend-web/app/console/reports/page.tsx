const mockReports = [
  { type: "Water has changed", place: "Lower Valley", submitted: "09:30", status: "Under review" },
  { type: "Flooding is visible", place: "East Ward", submitted: "08:15", status: "Unreviewed" },
  { type: "Conditions are unusually dry", place: "North Basin", submitted: "Yesterday", status: "Verified" },
];

export default function ConsoleReports() {
  return (
    <div className="space-y-6">
      <h1 className="text-2xl font-bold text-text-charcoal">Community reports</h1>

      {/* Tabs */}
      <div className="flex gap-2 border-b border-background-stone">
        {["Unreviewed", "Under review", "Verified", "Resolved"].map((tab, i) => (
          <button key={tab} className={`px-4 py-2 text-sm font-medium border-b-2 transition-colors ${
            i === 1 ? "border-forest text-forest" : "border-transparent text-text-muted hover:text-text-charcoal"
          }`}>
            {tab}
          </button>
        ))}
      </div>

      {/* Table */}
      <div className="bg-background rounded-xl border border-background-stone overflow-hidden">
        <table className="w-full text-left text-sm">
          <thead className="bg-background-mist border-b border-background-stone text-text-muted">
            <tr>
              <th className="p-4 font-medium">Type</th>
              <th className="p-4 font-medium">Place</th>
              <th className="p-4 font-medium">Submitted</th>
              <th className="p-4 font-medium">Status</th>
              <th className="p-4 font-medium text-right">Actions</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-background-stone">
            {mockReports.map((report, i) => (
              <tr key={i} className="hover:bg-background-mist/50">
                <td className="p-4 text-text-charcoal font-medium">{report.type}</td>
                <td className="p-4 text-text-muted">{report.place}</td>
                <td className="p-4 text-text-muted">{report.submitted}</td>
                <td className="p-4">
                  <span className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${
                    report.status === "Verified" ? "bg-status-normal/10 text-status-normal" :
                    report.status === "Under review" ? "bg-status-watch/10 text-status-watch" :
                    "bg-background-stone text-text-muted"
                  }`}>
                    {report.status}
                  </span>
                </td>
                <td className="p-4 text-right">
                  <button className="text-forest hover:underline text-xs font-medium">Review</button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
