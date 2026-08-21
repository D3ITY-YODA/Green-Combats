package com.greencompass.data.mock

import com.greencompass.domain.model.*

object MockData {
    val places = listOf(
        UserPlace("1", "Lower Valley", true),
        UserPlace("2", "East Ward", false),
        UserPlace("3", "North Basin", false)
    )

    fun getTodayData(placeId: String): TodayData {
        return when (placeId) {
            "1" -> TodayData(
                placeName = "Lower Valley",
                status = TodayStatus("1 important update", "Flood warning: Heavy rainfall may cause flooding within 24 hours.", StatusType.IMPORTANT, "10:00"),
                updates = listOf(
                    PublicUpdate("u1", "Flood warning", "Heavy rainfall may cause flooding within 24 hours.", "Lower Valley", "10:00", "Local weather authority", true),
                    PublicUpdate("u2", "Water outlook", "Water availability may decline over the next two weeks.", "Lower Valley", "today", "Water Resources Authority", false)
                ),
                exploreSections = listOf(
                    ExploreSection("local", "Local outlook", "Conditions for the coming days"),
                    ExploreSection("water", "Water outlook", "Information about nearby water conditions"),
                    ExploreSection("seasonal", "Seasonal information", "Changes that may affect your area"),
                    ExploreSection("community", "Community updates", "Information shared by people nearby")
                )
            )
            "2" -> TodayData(
                placeName = "East Ward",
                status = TodayStatus("No important updates", "for your area", StatusType.NORMAL, "10:00"),
                updates = emptyList(),
                exploreSections = listOf(
                    ExploreSection("local", "Local outlook", "Conditions for the coming days"),
                    ExploreSection("seasonal", "Seasonal information", "Changes that may affect your area")
                )
            )
            "3" -> TodayData(
                placeName = "North Basin",
                status = TodayStatus("Information delayed", "The latest update for North Basin is not available yet. Last reliable update: 08:00", StatusType.DELAYED, "08:00"),
                updates = emptyList(),
                exploreSections = listOf(
                    ExploreSection("seasonal", "Seasonal information", "Changes that may affect your area")
                )
            )
            else -> getTodayData("1")
        }
    }
}
