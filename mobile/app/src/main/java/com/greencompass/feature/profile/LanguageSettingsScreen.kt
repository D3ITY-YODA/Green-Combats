package com.greencompass.feature.profile

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.Check
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import com.greencompass.core.ui.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun LanguageSettingsScreen(onBack: () -> Unit) {
    var selectedLanguage by remember { mutableStateOf("English") }
    val languages = listOf("English", "Kiswahili", "French", "Arabic")

    GreenCompassScaffold(
        title = "Language",
        navigationIcon = {
            IconButton(onClick = onBack) {
                Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back", tint = GreenCompassColors.Charcoal)
            }
        }
    ) { paddingValues ->
        LazyColumn(
            modifier = Modifier.fillMaxSize().padding(paddingValues).padding(horizontal = AppSpacing.lg),
            verticalArrangement = Arrangement.spacedBy(AppSpacing.sm)
        ) {
            items(languages) { language ->
                val isSelected = selectedLanguage == language
                Surface(
                    modifier = Modifier.fillMaxWidth().clickable { selectedLanguage = language },
                    shape = RoundedCornerShape(12.dp),
                    color = if (isSelected) GreenCompassColors.SoftSage else Color.White,
                    border = BorderStroke(1.dp, if (isSelected) GreenCompassColors.ForestGreen else GreenCompassColors.Stone)
                ) {
                    Row(
                        modifier = Modifier.fillMaxWidth().padding(AppSpacing.md),
                        horizontalArrangement = Arrangement.SpaceBetween,
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Text(text = language, style = GreenCompassTypography.bodyLarge, color = GreenCompassColors.Charcoal)
                        if (isSelected) {
                            Icon(Icons.Default.Check, contentDescription = "Selected", tint = GreenCompassColors.ForestGreen)
                        }
                    }
                }
            }
        }
    }
}
