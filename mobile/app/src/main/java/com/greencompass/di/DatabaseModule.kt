package com.greencompass.di

import android.content.Context
import androidx.room.Room
import com.greencompass.data.local.GreenCompassDatabase
import com.greencompass.data.local.dao.ReportDao
import com.greencompass.data.local.dao.UpdateDao
import dagger.Module
import dagger.Provides
import dagger.hilt.InstallIn
import dagger.hilt.android.qualifiers.ApplicationContext
import dagger.hilt.components.SingletonComponent
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
object DatabaseModule {

    @Provides
    @Singleton
    fun provideDatabase(@ApplicationContext context: Context): GreenCompassDatabase {
        return Room.databaseBuilder(
            context,
            GreenCompassDatabase::class.java,
            "green_compass_db"
        ).fallbackToDestructiveMigration() // Added for version bump
        .build()
    }

    @Provides @Singleton fun provideUpdateDao(database: GreenCompassDatabase) = database.updateDao()
    @Provides @Singleton fun provideReportDao(database: GreenCompassDatabase) = database.reportDao()
}
