package com.greencompass.feature.today

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.greencompass.data.repository.TodayRepository
import com.greencompass.domain.model.Place
import com.greencompass.domain.model.TodayData
import com.greencompass.domain.model.TodayStatus
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.update
import kotlinx.coroutines.launch
import javax.inject.Inject

data class TodayUiState(
    val isLoading: Boolean = true,
    val data: TodayData? = null,
    val isRefreshing: Boolean = false,
    val errorMessage: String? = null
)

@HiltViewModel
class TodayViewModel @Inject constructor(
    private val repository: TodayRepository
) : ViewModel() {

    private val currentPlaceId = "place_123" 

    private val _uiState = MutableStateFlow(TodayUiState())
    val uiState: StateFlow<TodayUiState> = _uiState.asStateFlow()

    init {
        // Simulate loading then show empty state (Backend will replace this)
        viewModelScope.launch {
            kotlinx.coroutines.delay(800) // Show skeleton for 800ms
            _uiState.update {
                it.copy(
                    isLoading = false,
                    data = TodayData(
                        place = Place(id = currentPlaceId, name = "Lower Valley"),
                        status = TodayStatus(title = "No important updates", message = "Everything looks normal today.", dataStatus = "current"),
                        updates = emptyList(), // Empty list triggers the "No important alerts" UI
                        sections = emptyList()
                    )
                )
            }
        }
    }

    fun refresh() {
        _uiState.update { it.copy(isRefreshing = true, errorMessage = null) }
        viewModelScope.launch {
            repository.refreshToday(currentPlaceId)
                .onFailure { error ->
                    _uiState.update { it.copy(isRefreshing = false, errorMessage = error.message) }
                }
        }
    }
}
