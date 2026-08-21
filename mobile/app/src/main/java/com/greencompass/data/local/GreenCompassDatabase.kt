package com.greencompass.data.local

import androidx.room.Database
import androidx.room.RoomDatabase
import com.greencompass.data.local.dao.ReportDao
import com.greencompass.data.local.dao.UpdateDao
import com.greencompass.data.local.entity.ReportEntity
import com.greencompass.data.local.entity.UpdateEntity

@Database(
    entities = [UpdateEntity::class, ReportEntity::class],
    version = 2,
    exportSchema = false
)
abstract class GreenCompassDatabase : RoomDatabase() {
    abstract fun updateDao(): UpdateDao
    abstract fun reportDao(): ReportDao
}
