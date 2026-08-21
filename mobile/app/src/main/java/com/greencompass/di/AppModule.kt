package com.greencompass.di

import com.greencompass.data.repository.GreenCompassRepository
import com.greencompass.data.repository.MockGreenCompassRepository
import dagger.Binds
import dagger.Module
import dagger.hilt.InstallIn
import dagger.hilt.components.SingletonComponent
import javax.inject.Singleton

@Module
@InstallIn(SingletonComponent::class)
abstract class AppModule {
    @Binds
    @Singleton
    abstract fun bindGreenCompassRepository(impl: MockGreenCompassRepository): GreenCompassRepository
}
