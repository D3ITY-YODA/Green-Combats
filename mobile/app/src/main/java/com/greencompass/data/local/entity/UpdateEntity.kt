package com.greencompass.data.local.entity

import androidx.room.Entity
import androidx.room.Index
import androidx.room.PrimaryKey

@Entity(
    tableName = "updates",
    indices = [
        Index(value = ["placeId"]),
        Index(value = ["updatedAt"])
    ]
)
data class UpdateEntity(
    @PrimaryKey val id: String,
    val placeId: String,
    val placeName: String?,
    val topicKey: String,
    val type: String,
    val priority: String,
    val title: String,
    val message: String,
    val validFrom: String,
    val validUntil: String?,
    val updatedAt: String,
    val sourceName: String?
)
