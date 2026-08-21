package com.greencompass.data.local.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import com.greencompass.data.local.entity.UpdateEntity
import kotlinx.coroutines.flow.Flow

@Dao
interface UpdateDao {

    @Query("SELECT * FROM updates WHERE placeId = :placeId ORDER BY updatedAt DESC")
    fun observeUpdates(placeId: String): Flow<List<UpdateEntity>>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun upsertAll(updates: List<UpdateEntity>)

    @Query("DELETE FROM updates WHERE placeId = :placeId")
    suspend fun clearUpdatesForPlace(placeId: String)
}
