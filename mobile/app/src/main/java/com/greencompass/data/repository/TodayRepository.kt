package com.greencompass.data.repository

import com.greencompass.data.local.dao.UpdateDao
import com.greencompass.data.mapper.toEntity
import com.greencompass.data.remote.GreenCompassApi
import com.greencompass.domain.model.Place
import com.greencompass.domain.model.PublicUpdate
import com.greencompass.domain.model.TodayData
import com.greencompass.domain.model.TodayStatus
import com.greencompass.domain.model.TopicSection
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map
import javax.inject.Inject
import javax.inject.Singleton

interface TodayRepository {
    fun observeToday(placeId: String): Flow<TodayData?>
    suspend fun refreshToday(placeId: String): Result<Unit>
}

@Singleton
class DefaultTodayRepository @Inject constructor(
    private val api: GreenCompassApi,
    private val updateDao: UpdateDao
) : TodayRepository {

    override fun observeToday(placeId: String): Flow<TodayData?> {
        // In a full app, we'd observe Place and Sections too. 
        // For now, we map the updates to a basic TodayData structure.
        return updateDao.observeUpdates(placeId).map { entities ->
            if (entities.isEmpty()) null
            else TodayData(
                place = Place(id = placeId, name = "Lower Valley"), // Hardcoded for now until Place DB is added
                status = TodayStatus(
                    title = "1 important update",
                    message = "There is an important update for your area.",
                    dataStatus = "current"
                ),
                updates = entities.map { it.toDomain() },
                sections = emptyList() // Hardcoded until Context API is added
            )
        }
    }

    override suspend fun refreshToday(placeId: String): Result<Unit> {
        return try {
            val response = api.getToday(placeId)
            if (response.data != null) {
                val updates = response.data.updates.map { it.toEntity() }
                updateDao.clearUpdatesForPlace(placeId)
                updateDao.upsertAll(updates)
            }
            Result.success(Unit)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    // Helper extension to map Entity to Domain
    private fun com.greencompass.data.local.entity.UpdateEntity.toDomain(): PublicUpdate {
        return PublicUpdate(
            id = this.id,
            title = this.title,
            message = this.message,
            placeName = this.placeName,
            updatedAt = this.updatedAt
        )
    }
}
