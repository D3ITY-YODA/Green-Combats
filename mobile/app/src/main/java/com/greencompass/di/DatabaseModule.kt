package com.greencompass.di

import android.content.Context
import androidx.room.Room
import com.greencompass.data.local.GreenCompassDatabase
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
        ).build()
    }

    @Provides
    @Singleton
    fun provideUpdateDao(database: GreenCompassDatabase) = database.updateDao()
}
