package com.greencompass.data.remote

import com.greencompass.data.remote.dto.ApiResponse
import com.greencompass.data.remote.dto.TodayResponseDto
import retrofit2.http.GET
import retrofit2.http.Header
import retrofit2.http.Path
import retrofit2.http.Query

interface GreenCompassApi {

    @GET("api/v1/places/{placeId}/today")
    suspend fun getToday(
        @Path("placeId") placeId: String,
        @Query("locale") locale: String = "en"
    ): ApiResponse<TodayResponseDto>

    // We will add /reports, /explore, /auth here as we build those features
}
