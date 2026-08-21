package com.greencompass.feature.today

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.greencompass.data.repository.GreenCompassRepository
import com.greencompass.domain.model.TodayData
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
    val errorMessage: String? = null
)

@HiltViewModel
class TodayViewModel @Inject constructor(
    private val repository: GreenCompassRepository
) : ViewModel() {

    private val currentPlaceId = "1" // "1" is Lower Valley in MockData

    private val _uiState = MutableStateFlow(TodayUiState())
    val uiState: StateFlow<TodayUiState> = _uiState.asStateFlow()

    init {
        observeData()
    }

    private fun observeData() {
        viewModelScope.launch {
            repository.observeToday(currentPlaceId).collect { data ->
                _uiState.update {
                    it.copy(
                        isLoading = false,
                        data = data
                    )
                }
            }
        }
    }
}
