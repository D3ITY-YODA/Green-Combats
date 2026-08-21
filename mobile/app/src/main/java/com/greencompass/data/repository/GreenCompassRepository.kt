package com.greencompass.data.repository

import com.greencompass.data.mock.MockData
import com.greencompass.domain.model.TodayData
import com.greencompass.domain.model.UserPlace
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.flowOf
import javax.inject.Inject
import javax.inject.Singleton

interface GreenCompassRepository {
    fun observeToday(placeId: String): Flow<TodayData>
    fun observeUserPlaces(): Flow<List<UserPlace>>
}

@Singleton
class MockGreenCompassRepository @Inject constructor() : GreenCompassRepository {
    override fun observeToday(placeId: String): Flow<TodayData> {
        return flowOf(MockData.getTodayData(placeId))
    }

    override fun observeUserPlaces(): Flow<List<UserPlace>> {
        return flowOf(MockData.places)
    }
}
