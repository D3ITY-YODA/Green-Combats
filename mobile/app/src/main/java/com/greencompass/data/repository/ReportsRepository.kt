package com.greencompass.data.repository

import com.greencompass.data.local.dao.ReportDao
import com.greencompass.data.local.entity.ReportEntity
import kotlinx.coroutines.flow.Flow
import java.util.UUID
import javax.inject.Inject
import javax.inject.Singleton

interface ReportsRepository {
    suspend fun submitReport(type: String, placeId: String, description: String?): Result<String>
    fun getPendingReports(): Flow<List<ReportEntity>>
}

@Singleton
class DefaultReportsRepository @Inject constructor(
    private val reportDao: ReportDao
) : ReportsRepository {

    override suspend fun submitReport(type: String, placeId: String, description: String?): Result<String> {
        return try {
            val reportId = UUID.randomUUID().toString()
            val entity = ReportEntity(
                id = reportId,
                type = type,
                placeId = placeId,
                description = description,
                observedAt = java.time.Instant.now().toString(),
                status = "pending_sync" // Will be synced by WorkManager later
            )
            reportDao.insertReport(entity)
            Result.success(reportId)
        } catch (e: Exception) {
            Result.failure(e)
        }
    }

    override fun getPendingReports(): Flow<List<ReportEntity>> = reportDao.getPendingReports()
}
