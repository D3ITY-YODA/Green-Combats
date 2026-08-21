package com.greencompass.data.local.entity

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "reports")
data class ReportEntity(
    @PrimaryKey val id: String,
    val type: String,
    val placeId: String,
    val description: String?,
    val observedAt: String,
    val status: String // "pending_sync", "submitted"
)
