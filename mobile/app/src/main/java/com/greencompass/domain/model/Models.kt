package com.greencompass.domain.model

data class TodayData(
    val placeName: String,
    val status: TodayStatus,
    val updates: List<PublicUpdate>,
    val exploreSections: List<ExploreSection>
)

data class TodayStatus(
    val title: String,
    val message: String,
    val type: StatusType,
    val updatedAt: String
)

enum class StatusType { NORMAL, IMPORTANT, DELAYED, OFFLINE }

data class PublicUpdate(
    val id: String,
    val title: String,
    val message: String,
    val placeName: String,
    val updatedAt: String,
    val source: String,
    val isImportant: Boolean
)

data class ExploreSection(
    val key: String,
    val title: String,
    val description: String
)

data class UserPlace(val id: String, val name: String, val isPrimary: Boolean)

data class Organization(val name: String, val type: String, val area: String)
