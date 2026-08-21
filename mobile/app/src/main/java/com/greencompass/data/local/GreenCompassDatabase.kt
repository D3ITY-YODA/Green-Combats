package com.greencompass.data.local

import androidx.room.Database
import androidx.room.RoomDatabase
import com.greencompass.data.local.dao.UpdateDao
import com.greencompass.data.local.entity.UpdateEntity

@Database(
    entities = [
        UpdateEntity::class
        // We will add PlaceEntity, ReportEntity here later
    ],
    version = 1,
    exportSchema = false
)
abstract class GreenCompassDatabase : RoomDatabase() {
    abstract fun updateDao(): UpdateDao
}
