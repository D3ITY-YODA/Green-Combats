package com.greencompass.feature.onboarding

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.Clear
import androidx.compose.material.icons.filled.Search
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SearchPlaceScreen(
    onBack: () -> Unit,
    onSelectPlace: (String) -> Unit
) {
    var query by remember { mutableStateOf("") }
    val recentPlaces = listOf("Lower Valley", "East Ward", "North Basin")

    GreenCompassScaffold(
        title = "",
        navigationIcon = {
            IconButton(onClick = onBack) {
                Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal)
            }
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(horizontal = AppSpacing.lg)
        ) {
            Text(
                text = "Search for a place",
                style = GreenCompassTypography.headlineLarge,
                color = GreenCompassColors.Charcoal,
                modifier = Modifier.padding(bottom = AppSpacing.xl)
            )

            OutlinedTextField(
                value = query,
                onValueChange = { query = it },
                label = { Text("Search places") },
                leadingIcon = { Icon(Icons.Default.Search, contentDescription = null) },
                trailingIcon = {
                    if (query.isNotEmpty()) {
                        IconButton(onClick = { query = "" }) {
                            Icon(Icons.Default.Clear, contentDescription = "Clear")
                        }
                    }
                },
                modifier = Modifier.fillMaxWidth().padding(bottom = AppSpacing.xl),
                shape = RoundedCornerShape(12.dp)
            )

            Text(
                text = "Recent places",
                style = GreenCompassTypography.labelMedium,
                color = GreenCompassColors.MutedText,
                modifier = Modifier.padding(bottom = AppSpacing.sm)
            )

            LazyColumn(verticalArrangement = Arrangement.spacedBy(AppSpacing.sm)) {
                items(recentPlaces) { place ->
                    Text(
                        text = place,
                        style = GreenCompassTypography.bodyLarge,
                        color = GreenCompassColors.Charcoal,
                        modifier = Modifier
                            .fillMaxWidth()
                            .clickable { onSelectPlace(place) }
                            .padding(vertical = AppSpacing.sm)
                    )
                    Divider(color = GreenCompassColors.Stone)
                }
            }
        }
    }
}
