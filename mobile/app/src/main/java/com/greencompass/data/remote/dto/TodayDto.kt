package com.greencompass.data.remote.dto

import com.google.gson.annotations.SerializedName

data class TodayResponseDto(
    @SerializedName("place") val place: PlaceDto,
    @SerializedName("status") val status: TodayStatusDto,
    @SerializedName("updates") val updates: List<UpdateDto>,
    @SerializedName("sections") val sections: List<TopicSectionDto>,
    @SerializedName("updated_at") val updatedAt: String
)

data class PlaceDto(
    @SerializedName("id") val id: String,
    @SerializedName("name") val name: String,
    @SerializedName("type") val type: String,
    @SerializedName("country_code") val countryCode: String
)

data class TodayStatusDto(
    @SerializedName("title") val title: String,
    @SerializedName("message") val message: String,
    @SerializedName("updated_at") val updatedAt: String? = null,
    @SerializedName("data_status") val dataStatus: String
)

data class UpdateDto(
    @SerializedName("id") val id: String,
    @SerializedName("place_id") val placeId: String,
    @SerializedName("place_name") val placeName: String? = null,
    @SerializedName("topic_key") val topicKey: String,
    @SerializedName("type") val type: String,
    @SerializedName("priority") val priority: String,
    @SerializedName("title") val title: String,
    @SerializedName("message") val message: String,
    @SerializedName("valid_from") val validFrom: String,
    @SerializedName("valid_until") val validUntil: String? = null,
    @SerializedName("updated_at") val updatedAt: String,
    @SerializedName("source_name") val sourceName: String? = null
)

data class TopicSectionDto(
    @SerializedName("key") val key: String,
    @SerializedName("title") val title: String,
    @SerializedName("description") val description: String,
    @SerializedName("available") val available: Boolean
)
