package com.greencompass.domain.model

data class TodayData(
    val place: Place,
    val status: TodayStatus,
    val updates: List<PublicUpdate>,
    val sections: List<TopicSection>
)

data class Place(
    val id: String,
    val name: String
)

data class TodayStatus(
    val title: String,
    val message: String,
    val dataStatus: String // "current", "delayed", "stale"
)

data class PublicUpdate(
    val id: String,
    val title: String,
    val message: String,
    val placeName: String?,
    val updatedAt: String
)

data class TopicSection(
    val key: String,
    val title: String,
    val description: String
)
